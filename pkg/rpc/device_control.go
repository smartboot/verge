package rpc

import (
	"github.com/ibuilding-x/driver-box/driverbox"
	"github.com/ibuilding-x/driver-box/driverbox/helper"

	"github.com/ibuilding-x/driver-box/driverbox/plugin"
	"go.uber.org/zap"
)

func HandleDeviceControl(ctx Context, params interface{}) error {
	helper.Logger.Info("Handling device control", zap.Any("params", params))

	// Define structure for device control parameters
	type DeviceControlParams struct {
		ID     string            `json:"id"`
		Points map[string]string `json:"points"`
	}

	var controlParams DeviceControlParams
	err := helper.Map2Struct(params, &controlParams)
	if err != nil {
		return err
	}
	pointData := make([]plugin.PointData, 0)
	for pointName, pointValue := range controlParams.Points {
		pointData = append(pointData, plugin.PointData{
			PointName: pointName,
			Value:     pointValue,
		})
	}
	return driverbox.WritePoints(controlParams.ID, pointData)
}
