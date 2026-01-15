# Verge Export

## Introduction

Verge Export is a Go-based export service project designed to report device data, product information, and device status to a remote server. This project implements JSON-RPC protocol communication and supports multiple functionalities including device control, device addition, configuration changes, and network status monitoring.

## Main Features

- **Device Management**: Supports device information reporting and status synchronization
- **Product Management**: Collects and reports product information with MD5 checksum verification
- **Real-time Communication**: Long-lived connection mechanism based on Server-Sent Events (SSE)
- **RPC Handlers**: Built-in multiple RPC handlers including device control, command execution, and configuration changes
- **Security Authentication**: Token-based authentication mechanism to ensure secure communication

## Project Structure

```
verge-export/
├── export.go              # Main export logic implementation
├── main.go                # Program entry point
├── model.go               # Data model definitions
├── pkg/
│   ├── reporter/          # Data reporting module
│   │   ├── devices.go     # Device information reporting
│   │   ├── http_client.go # HTTP client wrapper
│   │   ├── products.go    # Product information reporting
│   │   ├── reporter.go    # Core reporting functionality
│   │   └── shadows.go     # Device status reporting
│   ├── rpc/               # JSON-RPC handlers
│   │   ├── context.go     # RPC context definition
│   │   ├── device_control.go  # Device control handling
│   │   ├── devices_add.go     # Device addition handling
│   │   ├── handlers.go        # RPC handler registration
│   │   ├── node_command.go    # Node command handling
│   │   ├── node_config_changed.go  # Configuration change handling
│   │   ├── node_network_status.go  # Network status handling
│   │   ├── product_import.go  # Product import handling
│   │   ├── products_report.go # Product report handling
│   │   └── types.go           # Type definitions
│   └── sse/                # SSE communication module
│       └── sse_manager.go  # SSE connection management
```

## Environment Requirements

- Go 1.16 or higher
- Linux/macOS/Windows operating systems

## Configuration

Configure service connection parameters via environment variables:

| Environment Variable | Description | Default Value |
|----------------------|-------------|---------------|
| `ENV_VERGE_BASE_URL` | Base URL of the remote server | - |

## Quick Start

### Build the Project

```bash
make
```

### Run the Service

```bash
./verge-export
```

## Core API

### Export Struct

Provides the main functional interfaces for the export service:

- `Init()` - Initialize the export service
- `Destroy()` - Destroy service resources
- `IsReady()` - Check service readiness status
- `GetBaseURL()` - Get the base URL
- `GetToken()` - Get the authentication token
- `ExportTo(deviceData)` - Export device data
- `OnEvent(eventCode, key, eventValue)` - Event handling

### Data Reporting

- `ReportDevices(deviceIds)` - Report device information
- `ReportShadows(deviceIds)` - Report device status
- `CollectAndReportProducts()` - Collect and report product information
- `ReportProducts(products)` - Directly report product information

## RPC Handlers

The project includes the following built-in RPC handlers:

- `HandleDeviceControl` - Device control commands
- `HandleDeviceAdd` - Add new devices
- `HandleCommand` - Execute node commands
- `HandleConfigChanged` - Configuration change notifications
- `HandleNetworkStatus` - Network status reporting
- `HandleProductImport` - Product resource import
- `HandleProductsReport` - Product information reporting

## License

This project is open-source. For specific licensing details, please refer to the project repository.