# Verge - 标准化边缘计算网关解决方案

[![CI](https://github.com/smartboot/verge/actions/workflows/ci.yml/badge.svg)](https://github.com/smartboot/verge/actions/workflows/ci.yml)
[![Docker](https://github.com/smartboot/verge/actions/workflows/docker.yml/badge.svg)](https://github.com/smartboot/verge/actions/workflows/docker.yml)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

> **Verge** 是一个轻量级、标准化的边缘计算网关解决方案。它将后台管理功能剥离为独立的云端/桌面服务，网关本身只专注于设备数据采集与治理，将硬件资源需求降至极致。

## 🎯 设计哲学

### 传统网关的困境

传统边缘网关往往集成设备接入、场景联动、规则引擎、用户界面等全量功能。这种"大而全"的设计带来了显著问题：

- **资源浪费** — 大量功能仅在实施阶段使用，却需长期占用网关资源
- **OTA 成本高** — 应用层频繁迭代导致升级包大、耗时长、风险高
- **硬件门槛高** — 功能臃肿推高了对 CPU、内存的最低要求

### Verge 的答案：关注点分离

```
┌─────────────────────────────────────┐
│      应用层（云端/桌面服务）          │
│   设备管理 · 场景联动 · 用户界面     │
│   ✓ 快速迭代  ✓ 丰富体验           │
└─────────────────┬───────────────────┘
                  │ 标准 API
┌─────────────────▼───────────────────┐
│      边缘层（Verge 网关）            │
│   设备接入 · 数据采集 · 协议转换     │
│   ✓ 极致轻量  ✓ 高度稳定           │
└─────────────────┬───────────────────┘
                  │ Modbus/MQTT/...
┌─────────────────▼───────────────────┐
│            物理设备                 │
└─────────────────────────────────────┘
```

**核心价值：**

| 维度 | 传统网关 | Verge 网关 |
|------|----------|------------|
| 资源占用 | 高 | 极致低 |
| OTA 频率 | 高 | 几乎为零 |
| 硬件要求 | 中高配置 | 任意硬件 |
| 稳定性 | 受应用迭代影响 | 长期稳定 |

---

## 🚀 快速开始

### 二进制部署（生产推荐）

从 [Releases](https://github.com/smartboot/verge/releases) 下载对应平台压缩包：

```bash
# 下载并解压
tar -xzf verge-linux-arm64-*.tar.gz
cd verge

# 配置并启动
export ENV_VERGE_BASE_URL=http://your-server:8080
./start.sh
```

**为什么推荐二进制部署？**
- 网关程序稳定后几乎无需更新，一次性部署即可
- 直接运行，无容器层开销
- 更适合工业现场环境

### Docker 部署（家用/开发场景）

适用于家庭实验室、开发测试环境，或设备基于以太网通信的场景：

```bash
docker-compose up -d
```

> **注意**：如果设备使用串口通信，需确保容器有相应权限访问串口设备。

访问 `http://localhost:18080` 查看管理界面。

### 源码运行（开发调试）

```bash
git clone https://github.com/smartboot/verge.git
cd verge && go mod tidy
export ENV_VERGE_BASE_URL=http://your-server:8080
go run cmd/main.go
```

---

## 📦 核心特性

- **极致轻量** — 聚焦设备接入核心功能，资源占用最小化
- **跨平台** — 支持 Linux/Windows/macOS，兼容 x86_64/ARM 架构
- **多协议** — 内置 MQTT、Modbus 等工业协议，支持 Lua 自定义扩展
- **云边协同** — SSE 实时通信 + JSON-RPC 标准接口
- **容器化** — 提供 Docker 镜像和 compose 配置

---

## 🔧 配置

### 环境变量

| 变量 | 说明 |
|------|------|
| `ENV_VERGE_BASE_URL` | 管理服务地址（必填）|

### 目录结构

```
verge/
├── cmd/           # 应用入口
├── pkg/           # 核心模块（SSE/RPC/Reporter）
├── res/           # 资源文件（驱动/模型/协议）
├── platform/      # 平台启动脚本
├── Dockerfile
└── docker-compose.yml
```

---

## 📡 架构说明

### 数据流

1. **启动** → 登录获取 token → 建立 SSE 连接
2. **上报** → 设备列表、产品模型、驱动信息
3. **接收** → 云端指令 → 执行设备控制/配置更新
4. **采集** → 设备数据 → 协议解析 → 定时上报
