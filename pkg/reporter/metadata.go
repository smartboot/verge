// Package reporter 提供数据上报功能
package reporter

import (
	"runtime"
	"time"

	"github.com/ibuilding-x/driver-box/v2/driverbox"
	"github.com/ibuilding-x/driver-box/v2/pkg/config"
	"github.com/smartboot/verge/pkg"
	"go.uber.org/zap"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

// appStartTime 记录应用程序启动时间
var appStartTime = time.Now()

// Metadata 节点元数据结构，包含基础的跨平台运行时信息
type Metadata struct {
	config.Metadata        // 嵌入driver-box的Metadata结构
	BuildTime       string `json:"buildTime"`
	Platform        string `json:"platform"`     // 运行平台（如：linux, darwin, windows）
	Architecture    string `json:"architecture"` // 系统架构（如：amd64, arm64）
	AppStartTime    int64  `json:"appStartTime"` // 应用程序启动时间戳
	AppUptime       int64  `json:"appUptime"`    // 应用运行时间(秒)
	Timestamp       int64  `json:"timestamp"`    // 上报时间戳
	SystemMemory    uint64 `json:"systemMemory"` // 系统内存(byte)
	AvailMemory     uint64 `json:"availMemory"`  // 剩余内存(byte)
	AppMemory       uint64 `json:"appMemory"`    // 应用内存(byte)
	// CPU负载信息
	CPUPercent float64 `json:"cpuPercent"` // CPU使用率(%)
	Load1      float64 `json:"load1"`      // 1分钟负载平均值
	Load5      float64 `json:"load5"`      // 5分钟负载平均值
	Load15     float64 `json:"load15"`     // 15分钟负载平均值
	NumCPU     int     `json:"numCPU"`     // CPU核心数
}

// ReportMetadata 上报节点元数据信息到服务器
// 收集当前节点的系统信息、运行状态等，并通过HTTP POST请求发送到服务器的/report/metadata端点
func (r *Reporter) ReportMetadata() error {
	driverbox.Log().Info("Reporting metadata")

	// 获取系统内存信息
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		driverbox.Log().Error("Failed to get system memory info", zap.Error(err))
		return err
	}

	// 获取基本CPU信息
	cpuCount := runtime.NumCPU()

	// 获取准确的CPU使用率
	cpuPercents, err := cpu.Percent(0, false)
	if err != nil {
		driverbox.Log().Warn("Failed to get CPU percent", zap.Error(err))
		cpuPercents = []float64{0}
	}

	// 获取准确的系统负载信息
	loadAvg, err := load.Avg()
	if err != nil {
		driverbox.Log().Warn("Failed to get load average", zap.Error(err))
	}

	// 获取当前进程内存信息
	psStat := &runtime.MemStats{}
	runtime.ReadMemStats(psStat)

	// 收集元数据信息
	metadata := &Metadata{
		Metadata:     driverbox.GetMetadata(), // 获取driver-box的基本元数据
		Platform:     runtime.GOOS,            // 运行平台
		Architecture: runtime.GOARCH,          // 系统架构
		BuildTime:    pkg.BuildTime,
		AppStartTime: appStartTime.Unix(),                       // 应用程序启动时间戳
		AppUptime:    int64(time.Since(appStartTime).Seconds()), // 应用运行时间(秒)
		Timestamp:    time.Now().Unix(),                         // 当前时间戳
		SystemMemory: vmStat.Total,                              // 系统内存(byte)
		AvailMemory:  vmStat.Available,
		AppMemory:    psStat.HeapSys, // 应用内存(byte)
		// CPU和负载信息
		CPUPercent: cpuPercents[0], // 准确的CPU使用率(%)
		Load1:      loadAvg.Load1,  // 准确的1分钟负载平均值
		Load5:      loadAvg.Load5,  // 准确的5分钟负载平均值
		Load15:     loadAvg.Load15, // 准确的15分钟负载平均值
		NumCPU:     cpuCount,       // CPU核心数
	}

	// 使用现有的postReport方法上报metadata
	err = r.postReport("report/metadata", metadata)
	if err != nil {
		driverbox.Log().Error("Failed to report metadata", zap.Error(err))
		return err
	}

	driverbox.Log().Info("Metadata report completed successfully")
	return nil
}
