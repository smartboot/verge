package reporter

import (
	"github.com/ibuilding-x/driver-box/v2/driverbox"
	"github.com/ibuilding-x/driver-box/v2/pkg/config"
	"go.uber.org/zap"
)

// ReportDevices sends device data to the server
func (r *Reporter) ReportDevices(deviceIds []string) error {
	type Model struct {
		config.Model
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
		device, ok := driverbox.CoreCache().GetDevice(deviceId)
		if !ok {
			driverbox.Log().Error("device not found", zap.String("deviceId", deviceId))
			continue
		}
		model, ok := driverbox.CoreCache().GetModel(device.ModelName)
		if !ok {
			driverbox.Log().Error("model not found", zap.String("modelName", device.ModelName))
			continue
		}
		_, connection := driverbox.CoreCache().GetConnection(device.ConnectionKey)
		if connection == nil {
			driverbox.Log().Error("connection not found", zap.String("connectionKey", device.ConnectionKey))
			continue
		}

		devices = append(devices, ReportDevice{
			Name:   device.Description,
			Device: device,
			Model: Model{
				Model:        model,
				DevicePoints: model.DevicePoints,
			},
			Connections:  connection,
			ProtocolName: device.PluginName,
		})
	}

	return r.postReport("report/devices", devices)
}
