package rpc

import (
	"github.com/ibuilding-x/driver-box/v2/driverbox"
	"go.uber.org/zap"
)

func HandleProductsReport(ctx Context, params interface{}) error {
	driverbox.Log().Info("Handling products report", zap.Any("params", params))

	// Collect and report products
	if err := ctx.CollectAndReportProducts(); err != nil {
		driverbox.Log().Error("Failed to collect and report products", zap.Error(err))
		return err
	}

	driverbox.Log().Info("Products report completed successfully")
	return nil
}
