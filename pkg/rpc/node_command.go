package rpc

import (
	"github.com/ibuilding-x/driver-box/pkg/driverbox/helper"
	"go.uber.org/zap"
)

func HandleCommand(ctx Context, params interface{}) error {
	helper.Logger.Info("Handling command", zap.Any("params", params))
	return nil
}
