// Package rpc 提供RPC相关的类型定义
package rpc

// RestResult REST API响应结果结构
type RestResult struct {
	Code    int         `json:"code"`    // 响应码
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 响应数据
}
