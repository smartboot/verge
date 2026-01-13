package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/ibuilding-x/driver-box/driverbox"
	"github.com/ibuilding-x/driver-box/driverbox/helper"
	"github.com/ibuilding-x/driver-box/driverbox/pkg/event"

	"github.com/ibuilding-x/driver-box/driverbox/plugin"
	"go.uber.org/zap"

	"github.com/smartboot/verge-export/pkg/reporter"
	"github.com/smartboot/verge-export/pkg/rpc"
	"github.com/smartboot/verge-export/pkg/sse"
)

var driverInstance *Export
var once = &sync.Once{}

const (
	ENV_VERGE_BASE_URL = "ENV_VERGE_BASE_URL"
)

// 设备自动发现插件
type Export struct {
	baseURL    string
	token      string
	ready      bool
	sseManager *sse.SSEManager
	reporter   *reporter.Reporter
}

func (export *Export) Init() error {
	export.ready = true
	helper.Crontab.AddFunc("10s", func() {
		if export.reporter == nil {
			return
		}
		deviceIds := make([]string, 0)
		for _, device := range helper.CoreCache.Devices() {
			deviceIds = append(deviceIds, device.ID)
		}
		export.ReportShadows(deviceIds)
	})
	return nil
}

// JSON-RPC 2.0 request structure
type JSONRPCRequest struct {
	Method  string      `json:"method"`
	JSONRPC string      `json:"jsonrpc"`
	Params  interface{} `json:"params"`
	ID      *int        `json:"id,omitempty"`
}

func (export *Export) handleJSONRPC(data string) error {
	var request JSONRPCRequest
	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return fmt.Errorf("failed to parse JSON-RPC request: %v", err)
	}

	// Verify JSON-RPC version
	if request.JSONRPC != "2.0" {
		return fmt.Errorf("unsupported JSON-RPC version: %s", request.JSONRPC)
	}

	// Handle different methods
	if handler, ok := rpc.Handlers[request.Method]; ok {
		return handler(export, request.Params)
	}
	helper.Logger.Warn("Unknown method", zap.String("method", request.Method))
	return nil
}

func (export *Export) Destroy() error {
	export.ready = false
	if export.sseManager != nil {
		export.sseManager.Disconnect()
		export.sseManager = nil
	}
	return nil
}
func NewExport() *Export {
	once.Do(func() {
		driverInstance = &Export{}
	})
	return driverInstance
}

// 点位变化触发场景联动
func (export *Export) ExportTo(deviceData plugin.DeviceData) {

}

// 继承Export OnEvent接口
func (export *Export) OnEvent(eventCode string, key string, eventValue interface{}) error {
	//网关启动完成
	if eventCode == event.EventCodeServiceStatus {
		export.login()
	}
	return nil
}

func (export *Export) IsReady() bool {
	return export.ready
}

func (export *Export) GetBaseURL() string {
	return export.baseURL
}

func (export *Export) GetToken() string {
	return export.token
}

func (export *Export) login() error {
	// Disconnect any existing SSE connection
	if export.sseManager != nil {
		export.sseManager.Disconnect()
	}

	// Get configuration from environment variables
	baseURL := os.Getenv(ENV_VERGE_BASE_URL)
	if baseURL == "" {
		return errors.New("VERGE_BASE_URL environment variable not set")
	}
	export.baseURL = strings.TrimSuffix(baseURL, "/")

	sn := driverbox.GetMetadata().SerialNo
	loginURL := export.baseURL + "/api/node/" + sn + "/login"
	// Prepare login payload
	loginData := map[string]string{"sn": sn}
	loginPayloadBytes, err := json.Marshal(loginData)
	if err != nil {
		return fmt.Errorf("failed to marshal login payload: %v", err)
	}
	loginPayload := string(loginPayloadBytes)

	// Login to get token
	resp, err := http.Post(loginURL, "application/json", strings.NewReader(loginPayload))
	if err != nil {
		return fmt.Errorf("failed to login: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode login response: %v", err)
	}

	token, ok := result["data"].(string)
	if !ok {
		return fmt.Errorf("token not found in response")
	}
	export.token = token

	// Create reporter
	export.reporter = reporter.NewReporter(export.baseURL, export.token)

	// Create and connect SSE manager
	export.sseManager = sse.NewSSEManager(export.baseURL, export.token, export.handleJSONRPC, func() {
		go func() {
			if err := export.login(); err != nil {
				helper.Logger.Error("Failed to re-login after token invalid", zap.Error(err))
			}
		}()
	})
	return export.sseManager.Connect()
}

// ReportDevices sends device data to the server
func (export *Export) ReportDevices(deviceIds []string) error {
	return export.reporter.ReportDevices(deviceIds)
}

func (export *Export) ReportShadows(deviceIds []string) error {
	return export.reporter.ReportShadows(deviceIds)
}

// CollectAndReportProducts collects product information from library and reports to server
func (export *Export) CollectAndReportProducts() error {
	return export.reporter.CollectAndReportProducts()
}

func (export *Export) ReportProducts(products []rpc.ProductInfo) error {
	return export.reporter.ReportProducts(products)
}
