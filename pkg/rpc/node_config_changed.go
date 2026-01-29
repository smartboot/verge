package rpc

import (
	"github.com/ibuilding-x/driver-box/v2/driverbox"
	"go.uber.org/zap"
)

func HandleConfigChanged(ctx Context, params interface{}) error {
	driverbox.Log().Info("Handling config change", zap.Any("params", params))
	return nil
}
