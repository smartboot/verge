// Package verge verge模块的核心类型定义
package verge

// RestResult REST API响应结果结构
type RestResult struct {
	Success bool        `json:"success"` // 是否成功
	Code    int         `json:"code"`    // 响应码
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 响应数据
}
