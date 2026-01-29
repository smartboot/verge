package rpc

import (
	"github.com/ibuilding-x/driver-box/v2/driverbox"
	"github.com/ibuilding-x/driver-box/v2/pkg/convutil"

	"github.com/ibuilding-x/driver-box/v2/driverbox/plugin"
	"go.uber.org/zap"
)

func HandleDeviceControl(ctx Context, params interface{}) error {
	driverbox.Log().Info("Handling device control", zap.Any("params", params))

	// Define structure for device control parameters
	type DeviceControlParams struct {
		ID     string            `json:"id"`
		Points map[string]string `json:"points"`
	}

	var controlParams DeviceControlParams
	err := convutil.Struct(params, &controlParams)
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
