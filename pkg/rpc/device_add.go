package rpc

import (
	"github.com/ibuilding-x/driver-box/pkg/driverbox/config"
	"github.com/ibuilding-x/driver-box/pkg/driverbox/helper"
	"go.uber.org/zap"
)

func HandleDeviceAdd(ctx Context, params interface{}) error {
	helper.Logger.Info("Handling device add", zap.Any("params", params))

	// Define structure for device add parameters
	type DeviceAddParams struct {
		Plugin        string             `json:"plugin"`
		Model         config.DeviceModel `json:"model"`
		ConnectionKey string             `json:"connectionKey"`
		Connection    any                `json:"connection"`
		Devices       []config.Device    `json:"devices"`
	}

	var addParams DeviceAddParams
	err := helper.Map2Struct(params, &addParams)
	if err != nil {
		return err
	}

	err = helper.CoreCache.AddModel(addParams.Plugin, addParams.Model)
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
