package rpc

var Handlers = map[string]func(Context, interface{}) error{
	"node.networkStatus": HandleNetworkStatus,
	"node.configChanged": HandleConfigChanged,
	"node.command":       HandleCommand,
	"device.control":     HandleDeviceControl,
	"devices.add":        HandleDeviceAdd,
	"product.import":     HandleProductImport,
	"products.report":    HandleProductsReport,
}
