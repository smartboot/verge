package rpc

var Handlers = map[string]func(Context, interface{}) error{
	"node.networkStatus": HandleNetworkStatus,
	"node.configChanged": HandleConfigChanged,
	"node.command":       HandleCommand,
	"device.control":     HandleDeviceControl,
	"devices.add":        HandleDeviceAdd,
	"devices.delete":     HandleDeviceDelete,
	"devices.report":     HandleDevicesReport, // 设备上报数据，未指定ID则全量上报
	"product.import":     HandleProductImport,
	"products.report":    HandleProductsReport,
}
