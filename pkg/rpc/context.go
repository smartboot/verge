package rpc

// ProductInfo 产品信息结构
type ProductInfo struct {
	Product string            `json:"product"`
	Hash    string            `json:"hash"`
	Models  map[string]string `json:"models"`
	Driver  map[string]string `json:"driver"`
}

// Context provides the interface for RPC handlers to interact with the main export functionality
type Context interface {
	ReportDevices(deviceIds []string) error
	ReportShadows(deviceIds []string) error
	ReportProducts(products []ProductInfo) error
	CollectAndReportProducts() error
	GetBaseURL() string
	GetToken() string
}
