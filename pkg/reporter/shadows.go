package reporter

import (
	"github.com/ibuilding-x/driver-box/v2/driverbox"
	"github.com/ibuilding-x/driver-box/v2/driverbox/shadow"
	"go.uber.org/zap"
)

func (r *Reporter) ReportShadows(deviceIds []string) error {
	driverbox.Log().Info("reporting shadows", zap.Int("deviceCount", len(deviceIds)))

	shadows := make([]shadow.Device, 0)
	for _, deviceId := range deviceIds {
		devShadow, ok := driverbox.Shadow().GetDevice(deviceId)
		if !ok {
			driverbox.Log().Error("shadow not found", zap.String("deviceId", deviceId))
			continue
		}
		shadows = append(shadows, devShadow)
	}

	return r.postReport("report/shadows", shadows)
}
