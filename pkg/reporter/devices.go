package reporter

import (
	"github.com/ibuilding-x/driver-box/pkg/driverbox/config"
	"github.com/ibuilding-x/driver-box/pkg/driverbox/helper"
	"go.uber.org/zap"
)

// ReportDevices sends device data to the server
func (r *Reporter) ReportDevices(deviceIds []string) error {
	type Model struct {
		config.ModelBase
		DevicePoints []config.Point `json:"devicePoints"`
	}
	type ReportDevice struct {
		//设备名称
		Name string `json:"name"`
		config.Device
		Model       Model `json:"model"`
		Connections any   `json:"connections" validate:""`
		// 协议名称（通过协议名称区分连接模式：客户端、服务端）
		ProtocolName string `json:"protocol" validate:"required"`
	}
	devices := make([]ReportDevice, 0)
	for _, deviceId := range deviceIds {
		device, ok := helper.CoreCache.GetDevice(deviceId)
		if !ok {
			helper.Logger.Error("device not found", zap.String("deviceId", deviceId))
			continue
		}
		model, ok := helper.CoreCache.GetModel(device.ModelName)
		if !ok {
			helper.Logger.Error("model not found", zap.String("modelName", device.ModelName))
			continue
		}
		connection, err := helper.CoreCache.GetConnection(device.ConnectionKey)
		if err != nil {
			helper.Logger.Error("connection not found", zap.String("connectionKey", device.ConnectionKey))
			continue
		}

		devices = append(devices, ReportDevice{
			Name:   device.Description,
			Device: device,
			Model: Model{
				ModelBase:    model.ModelBase,
				DevicePoints: model.DevicePoints,
			},
			Connections:  connection,
			ProtocolName: helper.CoreCache.GetConnectionPluginName(device.ConnectionKey),
		})
	}

	return r.postReport("report/devices", devices)
}
