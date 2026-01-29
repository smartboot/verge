// Package rpc 提供RPC上下文和类型定义
package rpc

// ProductInfo 产品信息结构，包含产品标识、哈希值、模型和驱动信息
type ProductInfo struct {
	Product string            `json:"product"` // 产品标识
	Hash    string            `json:"hash"`    // 产品哈希值
	Models  map[string]string `json:"models"`  // 模型映射 (modelKey -> hash)
	Driver  map[string]string `json:"driver"`  // 驱动映射 (driverKey -> hash)
}

// Context RPC处理器上下文接口，提供与主导出功能交互的方法
type Context interface {
	ReportDevices(deviceIds []string) error      // 上报设备数据
	ReportShadows(deviceIds []string) error      // 上报设备影子数据
	ReportProducts(products []ProductInfo) error // 上报产品信息
	CollectAndReportProducts() error             // 收集并上报所有产品
	GetBaseURL() string                          // 获取基础URL
	GetToken() string                            // 获取认证令牌
}
