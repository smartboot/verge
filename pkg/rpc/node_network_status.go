package rpc

import (
	"github.com/ibuilding-x/driver-box/v2/driverbox"
	"github.com/ibuilding-x/driver-box/v2/pkg/convutil"
	"go.uber.org/zap"
)

func HandleNetworkStatus(ctx Context, params interface{}) error {
	driverbox.Log().Info("Handling network status", zap.Any("params", params))
	//定义 networked 的结构体
	type NetworkStatus struct {
		Networked bool `json:"networked"`
	}
	networkStatus := NetworkStatus{}
	err := convutil.Struct(params, &networkStatus)
	if err != nil {
		return err
	}
	//组网成功，上报设备列表、模型和驱动文件列表
	if networkStatus.Networked {
		deviceIds := make([]string, 0)
		for _, device := range driverbox.CoreCache().Devices() {
			deviceIds = append(deviceIds, device.ID)
		}
		ctx.ReportDevices(deviceIds)
		ctx.ReportShadows(deviceIds)

		// Report products
		if err := ctx.CollectAndReportProducts(); err != nil {
			driverbox.Log().Error("Failed to report products", zap.Error(err))
			return err
		}
		driverbox.Log().Info("Networked, reporting device, model and driver lists")
	}
	return nil
}
