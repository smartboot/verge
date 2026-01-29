// Package main verge应用入口
package main

import (
	"os"

	"github.com/ibuilding-x/driver-box/v2/driverbox"
	"github.com/ibuilding-x/driver-box/v2/exports"
	"github.com/ibuilding-x/driver-box/v2/plugins"
	"github.com/smartboot/verge"
)

func main() {
	// 设置verge服务器基础URL环境变量
	os.Setenv(verge.ENV_VERGE_BASE_URL, "http://localhost:8080")

	// 注册所有插件
	plugins.EnableAll()

	// 加载所有导出器
	exports.EnableAll()

	// 加载verge导出器实例
	driverbox.EnableExport(verge.NewExport())

	// 启动driver-box框架
	driverbox.Start()

	// 阻塞主goroutine，保持程序运行
	select {}
}
