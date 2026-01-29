package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ibuilding-x/driver-box/v2/driverbox"
	"github.com/ibuilding-x/driver-box/v2/pkg/config"
	"github.com/ibuilding-x/driver-box/v2/pkg/convutil"
	"go.uber.org/zap"
)

func HandleProductImport(ctx Context, params interface{}) error {
	driverbox.Log().Info("Handling product import", zap.Any("params", params))

	// params should be an array of strings (resource paths)
	if params == nil {
		driverbox.Log().Error("Product import params is nil")
		return fmt.Errorf("product import params is nil")
	}

	// Convert params to []string
	var resourcePaths []string
	err := convutil.Struct(params, &resourcePaths)
	if err != nil {
		return err
	}

	// Process each resource path
	for _, resourcePath := range resourcePaths {
		driverbox.Log().Info("Processing resource path", zap.String("path", resourcePath))
		if err := importResource(ctx, resourcePath); err != nil {
			driverbox.Log().Error("Failed to import resource", zap.String("path", resourcePath), zap.Error(err))
			return fmt.Errorf("failed to import resource %s: %v", resourcePath, err)
		}
	}

	driverbox.Log().Info("Product import completed successfully", zap.Any("resourcePaths", resourcePaths))

	// Report products after import
	if err := ctx.CollectAndReportProducts(); err != nil {
		driverbox.Log().Error("Failed to report products after import", zap.Error(err))
		return err
	}

	return nil
}

func importResource(ctx Context, resourcePath string) error {
	fullURL := ctx.GetBaseURL() + resourcePath

	// Make HTTP request to fetch the resource
	resp, err := http.Get(fullURL)
	if err != nil {
		return fmt.Errorf("failed to fetch resource from %s: %v", fullURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch resource from %s, status: %d", fullURL, resp.StatusCode)
	}

	// Process the resource data based on content type
	contentType := resp.Header.Get("Content-Type")
	driverbox.Log().Info("Processing resource", zap.String("path", resourcePath), zap.String("contentType", contentType))

	if !strings.Contains(contentType, "application/json") {
		return errors.New("Content-Type is not application/json")
	}
	driverbox.Log().Info("Processing JSON resource", zap.String("path", resourcePath))

	// Parse the JSON to determine resource type
	var result RestResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode report devices response: %v", err)
	}

	if result.Code != 200 {
		return fmt.Errorf("report devices failed with code %d: %s", result.Code, result.Message)
	}
	type ModelResource struct {
		Name  string `json:"name"`
		Model string `json:"model"`
		Lua   string `json:"lua"`
	}
	type Resource struct {
		ProtocolKey string          `json:"protocolKey"`
		Lua         string          `json:"lua"`
		Models      []ModelResource `json:"models"`
	}

	res := new(Resource)
	err = convutil.Struct(result.Data, &res)
	if err != nil {
		return err
	}

	if len(res.ProtocolKey) > 0 {
		protocolDir := filepath.Join(config.ResourcePath, "library", "protocol")
		protocolPath := filepath.Join(protocolDir, res.ProtocolKey+".lua")

		if err := os.MkdirAll(protocolDir, 0755); err != nil {
			driverbox.Log().Error("Failed to create model directory", zap.String("dir", protocolDir), zap.Error(err))
			return fmt.Errorf("failed to create model directory: %v", err)
		}

		// Write model content to file
		if err := os.WriteFile(protocolPath, []byte(res.Lua), 0644); err != nil {
			driverbox.Log().Error("Failed to write protocol file", zap.String("path", protocolPath), zap.Error(err))
			return fmt.Errorf("failed to write protocol file: %v", err)
		}
		driverbox.Log().Info("protocol saved successfully", zap.String("path", protocolPath))
	}

	// Process each resource
	for _, resource := range res.Models {
		if resource.Name == "" {
			driverbox.Log().Error("Resource name is empty, skipping")
			continue
		}

		// Save model to resPath/library/model/name.json if model exists
		if resource.Model != "" {
			modelDir := filepath.Join(config.ResourcePath, "library", "model")
			modelPath := filepath.Join(modelDir, resource.Name+".json")

			// Create directory if it doesn't exist
			if err := os.MkdirAll(modelDir, 0755); err != nil {
				driverbox.Log().Error("Failed to create model directory", zap.String("dir", modelDir), zap.Error(err))
				return fmt.Errorf("failed to create model directory: %v", err)
			}

			// Write model content to file
			if err := os.WriteFile(modelPath, []byte(resource.Model), 0644); err != nil {
				driverbox.Log().Error("Failed to write model file", zap.String("path", modelPath), zap.Error(err))
				return fmt.Errorf("failed to write model file: %v", err)
			}
			driverbox.Log().Info("Model saved successfully", zap.String("path", modelPath))
		}

		// Save lua to resPath/library/driver/name.lua if lua exists
		if resource.Lua != "" {
			driverDir := filepath.Join(config.ResourcePath, "library", "driver")
			driverPath := filepath.Join(driverDir, resource.Name+".lua")

			// Create directory if it doesn't exist
			if err := os.MkdirAll(driverDir, 0755); err != nil {
				driverbox.Log().Error("Failed to create driver directory", zap.String("dir", driverDir), zap.Error(err))
				return fmt.Errorf("failed to create driver directory: %v", err)
			}

			// Write lua content to file
			if err := os.WriteFile(driverPath, []byte(resource.Lua), 0644); err != nil {
				driverbox.Log().Error("Failed to write lua file", zap.String("path", driverPath), zap.Error(err))
				return fmt.Errorf("failed to write lua file: %v", err)
			}
			driverbox.Log().Info("Lua file saved successfully", zap.String("path", driverPath))
		}
	}

	// For now, just log that we've processed the resource
	driverbox.Log().Info("JSON resource processed", zap.String("path", resourcePath))
	return nil
}
