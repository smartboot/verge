package reporter

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ibuilding-x/driver-box/driverbox"
	"github.com/ibuilding-x/driver-box/driverbox/helper"

	"go.uber.org/zap"

	"github.com/smartboot/verge-export/pkg/rpc"
)

// postReport performs a POST request to report data to the server
func (r *Reporter) postReport(endpoint string, payload interface{}) error {
	if !r.ready {
		return errors.New("reporter not ready")
	}
	// Get node serial number
	sn := driverbox.GetMetadata().SerialNo
	url := fmt.Sprintf("%s/api/node/%s/%s", r.baseURL, sn, endpoint)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal %s payload: %v", endpoint, err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return fmt.Errorf("failed to create %s request: %v", endpoint, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.token)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to %s: %v", endpoint, err)
	}
	defer resp.Body.Close()

	// Decode response
	var result rpc.RestResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode %s response: %v", endpoint, err)
	}

	if result.Code != 200 {
		return fmt.Errorf("%s failed with code %d: %s", endpoint, result.Code, result.Message)
	}

	helper.Logger.Info("Report successful", zap.String("endpoint", endpoint))
	return nil
}
