

# Verge Export

## 简介

Verge Export 是一个基于 Go 语言开发的导出服务项目，用于将设备数据、产品信息和设备状态报告到远程服务器。该项目实现了 JSON-RPC 协议通信，支持设备控制、设备添加、配置变更、网络状态监控等多种功能。

## 主要功能

- **设备管理**: 支持设备信息上报、设备状态同步
- **产品管理**: 产品信息收集与上报，支持 MD5 校验
- **实时通信**: 基于 Server-Sent Events (SSE) 的长连接通信机制
- **RPC 处理器**: 内置多种 RPC 处理器，包括设备控制、命令执行、配置变更等
- **安全认证**: Token 认证机制，确保通信安全

## 项目结构

```
verge-export/
├── export.go              # 主导出逻辑实现
├── main.go                # 程序入口
├── model.go               # 数据模型定义
├── pkg/
│   ├── reporter/          # 数据上报模块
│   │   ├── devices.go     # 设备信息上报
│   │   ├── http_client.go # HTTP 客户端封装
│   │   ├── products.go    # 产品信息上报
│   │   ├── reporter.go    # 报表核心功能
│   │   └── shadows.go     # 设备状态上报
│   ├── rpc/               # JSON-RPC 处理器
│   │   ├── context.go     # RPC 上下文定义
│   │   ├── device_control.go  # 设备控制处理
│   │   ├── devices_add.go     # 设备添加处理
│   │   ├── handlers.go        # RPC 处理器注册
│   │   ├── node_command.go    # 节点命令处理
│   │   ├── node_config_changed.go  # 配置变更处理
│   │   ├── node_network_status.go  # 网络状态处理
│   │   ├── product_import.go  # 产品导入处理
│   │   ├── products_report.go # 产品报告处理
│   │   └── types.go           # 类型定义
│   └── sse/                # SSE 通信模块
│       └── sse_manager.go  # SSE 连接管理
```

## 环境要求

- Go 1.16 或更高版本
- Linux/macOS/Windows 操作系统

## 配置说明

通过环境变量配置服务连接参数：

| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| `ENV_VERGE_BASE_URL` | 远程服务器基础地址 | - |

## 快速开始

### 编译项目

```bash
make
```

### 运行服务

```bash
./verge-export
```

## 核心 API

### Export 结构体

提供导出服务的主要功能接口：

- `Init()` - 初始化导出服务
- `Destroy()` - 销毁服务资源
- `IsReady()` - 检查服务就绪状态
- `GetBaseURL()` - 获取基础 URL
- `GetToken()` - 获取认证令牌
- `ExportTo(deviceData)` - 导出设备数据
- `OnEvent(eventCode, key, eventValue)` - 事件处理

### 数据上报

- `ReportDevices(deviceIds)` - 上报设备信息
- `ReportShadows(deviceIds)` - 上报设备状态
- `CollectAndReportProducts()` - 收集并上报产品信息
- `ReportProducts(products)` - 直接上报产品信息

## RPC 处理器

项目内置以下 RPC 处理器：

- `HandleDeviceControl` - 设备控制指令
- `HandleDeviceAdd` - 添加新设备
- `HandleCommand` - 执行节点命令
- `HandleConfigChanged` - 配置变更通知
- `HandleNetworkStatus` - 网络状态上报
- `HandleProductImport` - 产品资源导入
- `HandleProductsReport` - 产品信息报告

## License

本项目遵循开源协议，具体许可信息请查看项目仓库。