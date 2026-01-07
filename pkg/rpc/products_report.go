package rpc

import (
	"github.com/ibuilding-x/driver-box/driverbox/helper"
	"go.uber.org/zap"
)

func HandleProductsReport(ctx Context, params interface{}) error {
	helper.Logger.Info("Handling products report", zap.Any("params", params))

	// Collect and report products
	if err := ctx.CollectAndReportProducts(); err != nil {
		helper.Logger.Error("Failed to collect and report products", zap.Error(err))
		return err
	}

	helper.Logger.Info("Products report completed successfully")
	return nil
}
