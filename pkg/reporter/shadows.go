package reporter

import (
	"github.com/ibuilding-x/driver-box/pkg/driverbox/helper"
	"github.com/ibuilding-x/driver-box/pkg/driverbox/pkg/shadow"
	"go.uber.org/zap"
)

func (r *Reporter) ReportShadows(deviceIds []string) error {
	helper.Logger.Info("reporting shadows", zap.Int("deviceCount", len(deviceIds)))

	shadows := make([]shadow.Device, 0)
	for _, deviceId := range deviceIds {
		devShadow, ok := helper.DeviceShadow.GetDevice(deviceId)
		if !ok {
			helper.Logger.Error("shadow not found", zap.String("deviceId", deviceId))
			continue
		}
		shadows = append(shadows, devShadow)
	}

	return r.postReport("report/shadows", shadows)
}
