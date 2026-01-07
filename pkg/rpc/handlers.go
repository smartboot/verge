package rpc

var Handlers = map[string]func(Context, interface{}) error{
	"node.networkStatus": HandleNetworkStatus,
	"node.configChanged": HandleConfigChanged,
	"node.command":       HandleCommand,
	"device.control":     HandleDeviceControl,
	"device.add":         HandleDeviceAdd,
	"product.import":     HandleProductImport,
}
