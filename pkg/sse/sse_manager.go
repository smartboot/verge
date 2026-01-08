package sse

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"

	"github.com/ibuilding-x/driver-box/pkg/driverbox/helper"
	"go.uber.org/zap"
)

// SSEManager handles Server-Sent Events connection and reconnection
type SSEManager struct {
	baseURL        string
	token          string
	stopCh         chan struct{}
	sseResp        *http.Response
	onData         func(data string) error // callback for handling received data
	onTokenInvalid func()                  // callback when token is invalid
}

// NewSSEManager creates a new SSEManager instance
func NewSSEManager(baseURL, token string, onData func(data string) error, onTokenInvalid func()) *SSEManager {
	return &SSEManager{
		baseURL:        baseURL,
		token:          token,
		onData:         onData,
		onTokenInvalid: onTokenInvalid,
		stopCh:         make(chan struct{}),
	}
}

// Connect establishes the SSE connection and starts listening
func (sm *SSEManager) Connect() error {
	// Close any existing connection
	sm.Disconnect()

	sseURL := strings.TrimSuffix(sm.baseURL, "/") + "/api/node/sse/" + sm.token
	req, err := http.NewRequest("GET", sseURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create SSE request: %v", err)
	}
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")

	client := &http.Client{}
	sseResp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to establish SSE connection: %v", err)
	}

	if sseResp.StatusCode != http.StatusOK {
		sseResp.Body.Close()
		return fmt.Errorf("SSE connection failed with status: %d", sseResp.StatusCode)
	}

	sm.sseResp = sseResp
	sm.stopCh = make(chan struct{})

	// Start listening in a goroutine
	go sm.listen()
	return nil
}

// Disconnect closes the SSE connection
func (sm *SSEManager) Disconnect() {
	if sm.stopCh != nil {
		close(sm.stopCh)
		sm.stopCh = nil
	}
	if sm.sseResp != nil {
		sm.sseResp.Body.Close()
		sm.sseResp = nil
	}
}

// listen continuously listens for SSE events
func (sm *SSEManager) listen() {
	scanner := bufio.NewScanner(sm.sseResp.Body)
	for {
		select {
		case <-sm.stopCh:
			sm.sseResp.Body.Close()
			return
		default:
			if scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "data:") {
					data := strings.TrimPrefix(line, "data:")
					data = strings.TrimSpace(data)
					helper.Logger.Info("Received SSE data", zap.String("data", data))
					if err := sm.onData(data); err != nil {
						helper.Logger.Error("Error handling SSE data", zap.Error(err))
					}
				}
			} else {
				// Connection closed - token invalid, trigger re-login
				helper.Logger.Warn("SSE connection closed, assuming token invalid, calling onTokenInvalid")
				sm.sseResp.Body.Close()
				if sm.onTokenInvalid != nil {
					go sm.onTokenInvalid()
				}
				return
			}
		}
	}
}
