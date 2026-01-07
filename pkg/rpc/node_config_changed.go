package rpc

import (
	"github.com/ibuilding-x/driver-box/driverbox/helper"
	"go.uber.org/zap"
)

func HandleConfigChanged(ctx Context, params interface{}) error {
	helper.Logger.Info("Handling config change", zap.Any("params", params))
	return nil
}
