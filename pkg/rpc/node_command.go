package rpc

import (
	"github.com/ibuilding-x/driver-box/v2/driverbox"
	"go.uber.org/zap"
)

func HandleCommand(ctx Context, params interface{}) error {
	driverbox.Log().Info("Handling command", zap.Any("params", params))
	return nil
}
