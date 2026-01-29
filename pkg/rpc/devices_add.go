package rpc

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ibuilding-x/driver-box/v2/driverbox"
	"github.com/ibuilding-x/driver-box/v2/pkg/config"
	"github.com/ibuilding-x/driver-box/v2/pkg/convutil"
	"github.com/ibuilding-x/driver-box/v2/pkg/library"
	"go.uber.org/zap"
)

func HandleDeviceAdd(ctx Context, params interface{}) error {
	driverbox.Log().Info("Handling device add", zap.Any("params", params))

	// Define structure for device add parameters
	type DeviceAddParams struct {
		Plugin        string          `json:"plugin"`
		ModelKey      string          `json:"modelKey"`
		ModelHash     string          `json:"modelHash"`
		ConnectionKey string          `json:"connectionKey"`
		Connection    any             `json:"connection"`
		Devices       []config.Device `json:"devices"`
	}

	var addParams DeviceAddParams
	err := convutil.Struct(params, &addParams)
	if err != nil {
		return err
	}

	// Build model file path
	if len(addParams.ModelKey) > 0 {
		modelPath := filepath.Join(config.ResourcePath, "library", "model", addParams.ModelKey+".json")

		// Read model file
		modelContent, err := os.ReadFile(modelPath)
		if err != nil {
			driverbox.Log().Error("Failed to read model file", zap.String("modelKey", addParams.ModelKey), zap.String("path", modelPath), zap.Error(err))
			return fmt.Errorf("failed to read model file: %v", err)
		}

		// Calculate MD5 hash
		hash := md5.Sum(modelContent)
		computedHash := hex.EncodeToString(hash[:])

		// Verify model hash
		if computedHash != addParams.ModelHash {
			driverbox.Log().Error("Model hash mismatch", zap.String("modelKey", addParams.ModelKey), zap.String("expected", addParams.ModelHash), zap.String("computed", computedHash))
			return fmt.Errorf("model hash mismatch for %s", addParams.ModelKey)
		}

		// Load model from library
		model, err := library.Model().LoadLibrary(addParams.ModelKey)
		if err != nil {
			driverbox.Log().Error("Failed to load model from library", zap.String("modelKey", addParams.ModelKey), zap.Error(err))
			return fmt.Errorf("failed to load model: %v", err)
		}
		model.Name = addParams.ModelKey + "_" + computedHash
		err = driverbox.CoreCache().AddModel(addParams.Plugin, model)
		if err != nil {
			return err
		}

		for _, device := range addParams.Devices {
			err = driverbox.CoreCache().AddOrUpdateDevice(device)
			if err != nil {
				driverbox.Log().Error("Failed to add or update device", zap.String("deviceId", device.ID), zap.Error(err))
			}
		}
	}

	err = driverbox.CoreCache().AddConnection(addParams.Plugin, addParams.ConnectionKey, addParams.Connection)
	if err != nil {
		return err
	}

	//driverbox.ReloadPlugins()
	//// Report the added device
	//if err := ctx.ReportDevices([]string{addParams.ID}); err != nil {
	//	driverbox.Log().Error("Failed to report added device", zap.String("deviceId", addParams.ID), zap.Error(err))
	//	return err
	//}

	//driverbox.Log().Info("Device added successfully", zap.String("deviceId", addParams.ID))
	return nil
}
