# Verge

[![CI](https://github.com/smartboot/verge/actions/workflows/ci.yml/badge.svg)](https://github.com/smartboot/verge/actions/workflows/ci.yml)
[![Docker](https://github.com/smartboot/verge/actions/workflows/docker.yml/badge.svg)](https://github.com/smartboot/verge/actions/workflows/docker.yml)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

[🇺🇸 English](README.md) | [🇨🇳 简体中文](README_zh.md)

> A lightweight, standardized edge computing gateway solution — Focused on device data acquisition and governance, minimizing hardware resource requirements.

## Design Philosophy

Traditional edge gateways often integrate full-stack functionality including device access, scene linkage, rule engines, and user interfaces. This "all-in-one" design leads to:

- High resource consumption, requiring powerful hardware
- Frequent OTA updates, where any feature iteration may affect core stability
- High hardware barriers, difficult to deploy on low-end devices

**Verge adopts an "Application-Edge Separation" architecture**: Complex business logic (scene linkage, rule engines, user interfaces) is placed in the application layer, while the edge gateway only handles core device access and data acquisition. The application layer can be deployed on cloud servers or run as a local desktop application.

```
┌─────────────────────────────────────────────────────────────┐
│        Application Layer (Cloud Server / Desktop App)         │
│                                                              │
│   ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│   │ Device   │  │  Scene   │  │  Rule    │  │   User   │   │
│   │ Management│  │ Linkage  │  │  Engine  │  │Interface │   │
│   └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
│                                                              │
│              ✓ Rapid Iteration  ✓ Rich Experience           │
└─────────────────────────────┬───────────────────────────────┘
                              │
                              │ SSE + JSON-RPC (Standard API)
                              │
┌─────────────────────────────▼───────────────────────────────┐
│                   Edge Layer (Verge Gateway)                 │
│                                                              │
│   ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│   │  Device  │  │   Data   │  │ Protocol │  │ Cloud-Edge│   │
│   │  Access  │  │Collection│  │Conversion│  │  Comm    │   │
│   └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
│                                                              │
│             ✓ Ultra-Lightweight  ✓ High Stability           │
└─────────────────────────────┬───────────────────────────────┘
                              │
                              │ Modbus / MQTT / HTTP
                              │
┌─────────────────────────────▼───────────────────────────────┐
│                      Physical Devices                        │
│                                                              │
│      Meters / HVAC / Lighting / Sensors / PV Inverters ...  │
└─────────────────────────────────────────────────────────────┘
```

### Core Advantages

| Dimension | Traditional Gateway | Verge Gateway |
|-----------|---------------------|---------------|
| Resource Usage | High | Ultra-Low |
| OTA Frequency | High | Near Zero |
| Hardware Requirements | Mid-High | Any Hardware |
| Stability | Affected by App Iteration | Long-term Stable |

## Technical Implementation

Verge is built on the [driver-box](https://github.com/ibuilding-x/driver-box) framework, leveraging its mature device access capabilities.

### Cloud-Edge Communication

- **SSE (Server-Sent Events)** - Maintains long-polling connection with cloud for real-time command reception
- **JSON-RPC 2.0** - Standardized remote procedure call protocol

### Communication Flow

```
┌─────────┐                                          ┌─────────┐
│  Verge  │                                          │  Cloud  │
└────┬────┘                                          └────┬────┘
     │                                                    │
     │  1. Login (Authenticate with Serial Number)        │
     │ ──────────────────────────────────────────────►    │
     │                                                    │
     │  2. Return Token                                   │
     │ ◄──────────────────────────────────────────────    │
     │                                                    │
     │  3. Establish SSE Connection                       │
     │ ──────────────────────────────────────────────►    │
     │                                                    │
     │  4. Periodically Report Device Shadow & Metadata   │
     │ ──────────────────────────────────────────────►    │
     │                                                    │
     │  5. Send Control Commands (JSON-RPC)               │
     │ ◄──────────────────────────────────────────────    │
     │                                                    │
     │  6. Execute Device Operations, Return Results      │
     │ ──────────────────────────────────────────────►    │
     │                                                    │
```

### Supported JSON-RPC Methods

| Method | Description |
|--------|-------------|
| `device.control` | Device control (read/write device points) |
| `devices.add` | Add devices |
| `devices.delete` | Delete devices |
| `devices.report` | Device data reporting |
| `products.report` | Product information reporting |
| `product.import` | Import product model |
| `node.configChanged` | Configuration change notification |
| `node.networkStatus` | Network status query |
| `node.command` | Execute custom command |

## Quick Start

### Requirements

- Go 1.23+ (for source compilation)
- Docker & Docker Compose (for containerized deployment)

### Option 1: Docker Compose

Suitable for home labs and development/testing environments:

```bash
# Clone the repository
git clone https://github.com/smartboot/verge.git
cd verge

# Start services
docker-compose up -d

# View logs
docker-compose logs -f verge
```

Management UI: `http://localhost:18080`

### Option 2: Binary Deployment (Recommended for Production)

**Why binary deployment?**

- Once stable, the gateway rarely needs updates — deploy once, run indefinitely
- Direct execution without container overhead
- Better suited for industrial environments

Download from [Releases](https://github.com/smartboot/verge/releases):

```bash
# Extract
tar -xzf verge-linux-arm64-*.tar.gz
cd verge

# Configure cloud endpoint
export ENV_VERGE_BASE_URL=http://your-server:8080

# Start
./start.sh
```

### Option 3: Run from Source

```bash
git clone https://github.com/smartboot/verge.git
cd verge

export ENV_VERGE_BASE_URL=http://your-server:8080

go mod tidy
go run cmd/main.go
```

## Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `ENV_VERGE_BASE_URL` | Cloud management service URL | Yes |

### Directory Structure

```
verge/
├── cmd/main.go              # Application entry
├── pkg/
│   ├── sse/                 # SSE connection management
│   ├── rpc/                 # JSON-RPC handlers
│   └── reporter/            # Data reporting module
├── res/
│   ├── driver/              # Driver configs (device models, connections)
│   └── library/
│       ├── driver/          # Lua driver scripts
│       ├── model/           # Product model definitions
│       └── protocol/        # Protocol parsing scripts
├── platform/                # Platform-specific startup scripts
├── Dockerfile
└── docker-compose.yml
```

## Extension Development

### Add RPC Methods

Register in [`pkg/rpc/handlers.go`](pkg/rpc/handlers.go):

```go
var Handlers = map[string]func(Context, interface{}) error{
    "your.method": HandleYourMethod,
}
```

### Custom Device Drivers

Add Lua scripts in `res/library/driver/`. Refer to [driver-box documentation](https://github.com/ibuilding-x/driver-box).

## Related Projects

- [driver-box](https://github.com/ibuilding-x/driver-box) - Device access framework
- [verge-ui](https://github.com/smartboot/verge-ui) - Web management interface

## License

[Apache License 2.0](LICENSE)
