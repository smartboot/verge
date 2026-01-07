package rpc

import (
	"github.com/ibuilding-x/driver-box/driverbox/config"
	"github.com/ibuilding-x/driver-box/driverbox/helper"
	"go.uber.org/zap"
)

func HandleDeviceAdd(ctx Context, params interface{}) error {
	helper.Logger.Info("Handling device add", zap.Any("params", params))

	// Define structure for device add parameters
	type DeviceAddParams struct {
		ID         string             `json:"id"`
		Plugin     string             `json:"plugin"`
		Model      config.DeviceModel `json:"model"`
		Connection any                `json:"connection"`
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

	// Report the added device
	if err := ctx.ReportDevices([]string{addParams.ID}); err != nil {
		helper.Logger.Error("Failed to report added device", zap.String("deviceId", addParams.ID), zap.Error(err))
		return err
	}

	helper.Logger.Info("Device added successfully", zap.String("deviceId", addParams.ID))
	return nil
}
