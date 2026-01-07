package rpc

import (
	"github.com/ibuilding-x/driver-box/driverbox/helper"
	"go.uber.org/zap"
)

func HandleCommand(ctx Context, params interface{}) error {
	helper.Logger.Info("Handling command", zap.Any("params", params))
	return nil
}
