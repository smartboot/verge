package rpc

import (
	"github.com/ibuilding-x/driver-box/driverbox/helper"
	"github.com/ibuilding-x/driver-box/driverbox/library"
	"github.com/ibuilding-x/driver-box/driverbox/pkg/config"
	"go.uber.org/zap"
)

func HandleDeviceAdd(ctx Context, params interface{}) error {
	helper.Logger.Info("Handling device add", zap.Any("params", params))

	// Define structure for device add parameters
	type DeviceAddParams struct {
		Plugin        string          `json:"plugin"`
		ModelKey      string          `json:"model"`
		ConnectionKey string          `json:"connectionKey"`
		Connection    any             `json:"connection"`
		Devices       []config.Device `json:"devices"`
	}

	var addParams DeviceAddParams
	err := helper.Map2Struct(params, &addParams)
	if err != nil {
		return err
	}
	model, _ := library.Model().LoadLibrary(addParams.ModelKey)
	err = helper.CoreCache.AddModel(addParams.Plugin, model)
	if err != nil {
		return err
	}

	err = helper.CoreCache.AddConnection(addParams.Plugin, addParams.ConnectionKey, addParams.Connection)
	if err != nil {
		return err
	}

	for _, device := range addParams.Devices {
		err = helper.CoreCache.AddOrUpdateDevice(device)
		if err != nil {
			helper.Logger.Error("Failed to add or update device", zap.String("deviceId", device.ID), zap.Error(err))
		}
	}

	//driverbox.ReloadPlugins()
	//// Report the added device
	//if err := ctx.ReportDevices([]string{addParams.ID}); err != nil {
	//	helper.Logger.Error("Failed to report added device", zap.String("deviceId", addParams.ID), zap.Error(err))
	//	return err
	//}

	//helper.Logger.Info("Device added successfully", zap.String("deviceId", addParams.ID))
	return nil
}
