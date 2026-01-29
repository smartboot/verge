package rpc

import (
	"github.com/ibuilding-x/driver-box/v2/driverbox"
	"github.com/ibuilding-x/driver-box/v2/pkg/convutil"
	"go.uber.org/zap"
)

func HandleDeviceDelete(ctx Context, params interface{}) error {
	driverbox.Log().Info("Handling device add", zap.Any("params", params))

	// Define structure for device add parameters
	ids := make([]string, 0)

	err := convutil.Struct(params, &ids)
	if err != nil {
		return err
	}

	err = driverbox.CoreCache().BatchRemoveDevice(ids)
	if err != nil {
		return err
	}
	driverbox.ReloadPlugins()
	return nil
}
