// Package rpc 提供RPC处理器实现
package rpc

import (
	"github.com/ibuilding-x/driver-box/v2/driverbox"
	"github.com/ibuilding-x/driver-box/v2/pkg/convutil"
	"go.uber.org/zap"
)

// HandleDevicesReport 处理设备上报请求
// 当params为nil或空时，上报所有设备；否则上报指定的设备列表
func HandleDevicesReport(ctx Context, params interface{}) error {
	driverbox.Log().Info("Handling devices report", zap.Any("params", params))

	// 初始化设备ID列表
	deviceIds := make([]string, 0)

	// 如果提供了参数，尝试解析设备ID列表
	if params != nil {
		err := convutil.Struct(params, &deviceIds)
		if err != nil {
			driverbox.Log().Error("Failed to convert params", zap.Error(err))
		}
	}

	// 如果没有提供参数或解析失败，收集所有设备ID进行全量上报
	if params == nil {
		for _, device := range driverbox.CoreCache().Devices() {
			deviceIds = append(deviceIds, device.ID)
		}
	}

	// 上报设备数据
	ctx.ReportDevices(deviceIds)
	driverbox.Log().Info("Devices report completed successfully", zap.Int("deviceCount", len(deviceIds)))

	// 上报设备影子数据
	ctx.ReportShadows(deviceIds)
	driverbox.Log().Info("Shadows report completed successfully", zap.Int("deviceCount", len(deviceIds)))

	return nil
}
