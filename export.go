package verge

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ibuilding-x/driver-box/v2/driverbox"
	"github.com/ibuilding-x/driver-box/v2/pkg/config"
	"github.com/ibuilding-x/driver-box/v2/pkg/event"
	"github.com/smartboot/verge/pkg"

	"github.com/ibuilding-x/driver-box/v2/driverbox/plugin"
	"go.uber.org/zap"

	"github.com/smartboot/verge/pkg/reporter"
	"github.com/smartboot/verge/pkg/rpc"
	"github.com/smartboot/verge/pkg/sse"
)

var driverInstance *Export
var once = &sync.Once{}

const (
	ENV_VERGE_BASE_URL = "ENV_VERGE_BASE_URL"
)

// 设备自动发现插件
type Export struct {
	baseURL    string
	token      string
	ready      bool
	sseManager *sse.SSEManager
	reporter   *reporter.Reporter
}

func (export *Export) Init() error {
	driverbox.UpdateMetadata(func(metadata *config.Metadata) {
		metadata.SoftwareVersion = pkg.Version
	})
	export.ready = true

	// 每10秒上报设备影子数据
	driverbox.AddFunc("10s", func() {
		if export.reporter == nil {
			return
		}
		deviceIds := make([]string, 0)
		for _, device := range driverbox.CoreCache().Devices() {
			deviceIds = append(deviceIds, device.ID)
		}
		export.ReportShadows(deviceIds)
	})

	// 每5分钟上报元数据信息
	driverbox.AddFunc("10s", func() {
		if export.reporter == nil {
			return
		}
		if err := export.ReportMetadata(); err != nil {
			driverbox.Log().Error("Failed to report metadata periodically", zap.Error(err))
		}
	})

	//确保目录都存在
	for _, subDir := range []string{"driver", "model", "protocol"} {
		dir := filepath.Join(config.ResourcePath, "library", subDir)
		if err := os.MkdirAll(dir, 0755); err != nil {
			driverbox.Log().Error("Failed to create directory", zap.String("dir", dir), zap.Error(err))
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}
	return nil
}

// JSONRPCRequest JSON-RPC 2.0 请求结构
type JSONRPCRequest struct {
	Method  string      `json:"method"`       // RPC方法名
	JSONRPC string      `json:"jsonrpc"`      // JSON-RPC版本
	Params  interface{} `json:"params"`       // 方法参数
	ID      *int        `json:"id,omitempty"` // 请求ID，可选
}

// handleJSONRPC 处理JSON-RPC请求，路由到对应的处理器
func (export *Export) handleJSONRPC(data string) error {
	var request JSONRPCRequest
	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return fmt.Errorf("failed to parse JSON-RPC request: %v", err)
	}

	// Verify JSON-RPC version
	if request.JSONRPC != "2.0" {
		return fmt.Errorf("unsupported JSON-RPC version: %s", request.JSONRPC)
	}

	// Handle different methods
	if handler, ok := rpc.Handlers[request.Method]; ok {
		return handler(export, request.Params)
	}
	driverbox.Log().Warn("Unknown method", zap.String("method", request.Method))
	return nil
}

func (export *Export) Destroy() error {
	export.ready = false
	if export.sseManager != nil {
		export.sseManager.Disconnect()
		export.sseManager = nil
	}
	return nil
}
func NewExport() *Export {
	once.Do(func() {
		driverInstance = &Export{}
	})
	return driverInstance
}

// 点位变化触发场景联动
func (export *Export) ExportTo(deviceData plugin.DeviceData) {

}

// 继承Export OnEvent接口
func (export *Export) OnEvent(eventCode event.EventCode, key string, eventValue interface{}) error {
	//网关启动完成
	if eventCode == event.ServiceStatus && eventValue == event.ServiceStatusHealthy {
		for {
			err := export.login()
			if err == nil {
				break
			}
			driverbox.Log().Error("Failed to login", zap.Error(err))
			time.Sleep(5 * time.Second)
		}

	}
	return nil
}

func (export *Export) IsReady() bool {
	return export.ready
}

func (export *Export) GetBaseURL() string {
	return export.baseURL
}

func (export *Export) GetToken() string {
	return export.token
}

// login 执行登录流程，获取认证令牌并建立SSE连接
func (export *Export) login() error {
	// Disconnect any existing SSE connection
	if export.sseManager != nil {
		export.sseManager.Disconnect()
	}

	// Get configuration from environment variables
	baseURL := os.Getenv(ENV_VERGE_BASE_URL)
	if baseURL == "" {
		return errors.New("VERGE_BASE_URL environment variable not set")
	}
	export.baseURL = strings.TrimSuffix(baseURL, "/")

	sn := driverbox.GetMetadata().SerialNo
	loginURL := export.baseURL + "/api/node/" + sn + "/login"
	// Prepare login payload
	loginData := map[string]string{"sn": sn}
	loginPayloadBytes, err := json.Marshal(loginData)
	if err != nil {
		return fmt.Errorf("failed to marshal login payload: %v", err)
	}
	loginPayload := string(loginPayloadBytes)

	// Login to get token
	resp, err := http.Post(loginURL, "application/json", strings.NewReader(loginPayload))
	if err != nil {
		return fmt.Errorf("failed to login: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode login response: %v", err)
	}

	token, ok := result["data"].(string)
	if !ok {
		return fmt.Errorf("token not found in response")
	}
	export.token = token

	// Create reporter
	export.reporter = reporter.NewReporter(export.baseURL, export.token)

	// Create and connect SSE manager
	export.sseManager = sse.NewSSEManager(export.baseURL, export.token, export.handleJSONRPC, func() {
		go func() {
			for {
				if err := export.login(); err == nil {
					return
				}
				driverbox.Log().Error("Failed to re-login after token invalid", zap.Error(err))
				time.Sleep(5 * time.Second)
			}

		}()
	})
	return export.sseManager.Connect()
}

// ReportDevices 上报设备数据到服务器
func (export *Export) ReportDevices(deviceIds []string) error {
	return export.reporter.ReportDevices(deviceIds)
}

// ReportShadows 上报设备影子数据到服务器
func (export *Export) ReportShadows(deviceIds []string) error {
	return export.reporter.ReportShadows(deviceIds)
}

// ReportMetadata 上报节点元数据信息
func (export *Export) ReportMetadata() error {
	return export.reporter.ReportMetadata()
}

// CollectAndReportProducts 从库中收集产品信息并上报到服务器
func (export *Export) CollectAndReportProducts() error {
	return export.reporter.CollectAndReportProducts()
}

// ReportProducts 上报产品信息到服务器
func (export *Export) ReportProducts(products []rpc.ProductInfo) error {
	return export.reporter.ReportProducts(products)
}
