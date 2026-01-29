# Verge - IoT边缘计算网关

Verge是一个基于[driver-box](https://github.com/ibuilding-x/driver-box)框架的IoT边缘计算网关项目，用于连接物理设备与云端服务，实现设备数据采集、控制命令转发、以及双向通信等功能。

## 功能特性

- **双向通信**: 基于Server-Sent Events(SSE)实现云端到边缘的实时指令下发
- **设备管理**: 支持动态添加、删除、更新设备信息
- **数据同步**: 自动上报设备状态、设备影子数据到云端
- **产品模型管理**: 支持物模型的导入和上报
- **协议适配**: 支持多种物联网通信协议（如MQTT等）
- **设备控制**: 接收云端指令对设备进行远程控制
- **网络状态监控**: 实时监控网络连接状态并上报

## 架构设计

```
云端服务 ←→ Verge网关 ←→ 物理设备
     ↑                    ↑
   SSE通信              协议适配
     ↓                    ↓
JSON-RPC指令          数据采集
```

### 核心模块

- **SSE管理器** (pkg/sse/sse_manager.go): 负责与云端建立长连接，接收下行指令
- **RPC处理器** (pkg/rpc/): 处理云端下发的各类指令
- **数据上报器** (pkg/reporter/): 负责向云端上报设备数据
- **导出器** (export.go): 整合各模块功能，作为driver-box插件入口

### RPC指令集

| 方法 | 描述 |
|------|------|
| `node.networkStatus` | 处理网络状态变化，网络连通时上报设备、模型和驱动 |
| `node.configChanged` | 处理配置变更通知 |
| `node.command` | 执行节点级命令 |
| `device.control` | 控制指定设备 |
| `devices.add` | 添加新设备 |
| `product.import` | 从云端导入产品模型和协议脚本 |
| `products.report` | 上报产品信息 |

## 环境要求

- Go 1.18+
- [driver-box](https://github.com/ibuilding-x/driver-box) 框架

## 快速开始

### 1. 环境配置

设置云端服务地址：
```bash
export ENV_VERGE_BASE_URL=http://your-cloud-server:8080
```

### 2. 运行项目

```bash
# 克隆项目
git clone https://github.com/smartboot/verge.git
cd verge

# 下载依赖
go mod tidy

# 运行
go run cmd/main.go
```

### 3. 部署构建

使用提供的部署脚本构建多平台二进制文件：
```bash
./deploy.sh
```

该脚本会生成适用于不同操作系统的压缩包，包括：
- Linux (arm64, amd64, arm)
- Windows (amd64, arm64)
- macOS (amd64, arm64)

## 使用说明

### 设备添加
通过云端下发`devices.add`指令可动态添加设备，参数包含：
- 插件类型
- 模型标识符及哈希值验证
- 连接配置
- 设备列表

### 设备控制
通过云端下发`device.control`指令控制设备，参数包含：
- 设备ID
- 点位映射（属性名-值对）

### 数据上报
网关定期（默认每10秒）上报设备影子数据，包括设备状态和属性值。

## 配置说明

### 环境变量

- `ENV_VERGE_BASE_URL`: 云端服务基础URL
- `ENV_RESOURCE_PATH`: 资源文件路径（可选，默认为 ./res）

### 资源目录结构

```
res/
├── library/
    ├── model/      # 物模型定义文件
    └── protocol/   # 协议适配脚本
```

## 开发指南

### 添加新的RPC处理器

1. 在pkg/rpc/目录下创建新的处理器文件
2. 在handlers.go中注册处理器函数
3. 实现相应的业务逻辑

### 添加资源文件

物模型文件应放置在res/library/model/目录下，协议脚本应放置在res/library/protocol/目录下。
此外，驱动脚本应放置在res/library/driver/目录下，这些会在产品导入时自动生成。

## 目录结构

```
verge/
├── cmd/                    # 应用入口
│   └── main.go
├── pkg/                    # 核心功能包
│   ├── reporter/           # 数据上报模块
│   ├── rpc/                # RPC处理模块
│   └── sse/                # SSE通信模块
├── res/                    # 资源文件
│   └── library/
    ├── model/          # 物模型定义
    └── protocol/       # 协议脚本
├── Makefile               # 构建脚本
├── deploy.sh              # 部署脚本
├── export.go              # 插件导出入口
├── model.go               # 数据模型定义
└── README.md              # 项目文档
```

## API接口

### 登录接口
- URL: `/api/node/{serial_no}/login`
- 方法: POST
- 用途: 获取访问令牌

### SSE接口
- URL: `/api/node/sse/{token}`
- 方法: GET
- 用途: 建立服务器发送事件连接，接收云端指令

## 安全性

- 使用序列号和令牌进行身份验证
- 模型文件通过MD5哈希验证完整性
- 支持令牌自动刷新机制

## 调试与日志

系统使用zap日志库记录运行信息，主要包括：
- 网络连接状态
- RPC调用详情
- 设备操作记录
- 错误和警告信息
