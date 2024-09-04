package information

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"math"
	"os"
	"os/user"
	"runtime"
	"time"

	"noah/client/app/entitie"
	"noah/client/app/service"

	"noah/client/app/utils/network"
)

type Service struct {
}

func NewService() service.Information {
	return &Service{}
}

func (i Service) LoadClientSpecs() (*entitie.Client, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	username, err := user.Current()
	if err != nil {
		return nil, err
	}
	macAddress, err := network.GetMacAddress()
	if err != nil {
		return nil, err
	}
	return &entitie.Client{
		Hostname:    hostname,
		Username:    username.Name,
		UserID:      username.Username,
		OSName:      runtime.GOOS,
		OSArch:      runtime.GOARCH,
		MacAddress:  macAddress,
		IPAddress:   network.GetLocalIP(),
		Port:        "",
		FetchedUnix: time.Now().UTC().Unix(),
	}, nil
}

// GetSystemInfo 获取系统的基本信息，包括CPU、内存、磁盘和网络信息。
func (i Service) GetSystemInfo() (*entitie.SystemInfo, error) {
	var sysInfo = &entitie.SystemInfo{}

	// 获取CPU信息
	_, err := cpu.Info()
	if err != nil {
		return nil, err
	}
	cpuUsage, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}
	sysInfo.CpuUsage = roundToTwoDecimals(cpuUsage[0])

	// 获取内存信息
	memStats, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	sysInfo.MemoryTotal = coverToGb(memStats.Total)
	sysInfo.MemoryUsed = coverToGb(memStats.Used)
	sysInfo.MemoryFree = coverToGb(memStats.Free)
	sysInfo.MemoryAvailable = coverToGb(memStats.Available)
	sysInfo.MemoryUsedPercent = roundToTwoDecimals(memStats.UsedPercent)

	// 获取磁盘信息
	usage, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}
	sysInfo.DiskTotal = coverToGb(usage.Total)
	sysInfo.DiskUsed = coverToGb(usage.Used)
	sysInfo.DiskFree = coverToGb(usage.Free)

	// 获取带宽信息
	bandwidth, err := net.IOCounters(true)
	if err != nil {
		return nil, fmt.Errorf("failed to get network bandwidth: %w", err)
	}
	var totalRxBytes, totalTxBytes uint64
	for _, stats := range bandwidth {
		totalRxBytes += stats.BytesRecv
		totalTxBytes += stats.BytesSent
	}

	// 计算总带宽
	rxBytesPerSec := float64(totalRxBytes) / 1024 / 1024
	txBytesPerSec := float64(totalTxBytes) / 1024 / 1024

	// 转换为Mbps
	sysInfo.BandwidthIn = roundToTwoDecimals(rxBytesPerSec)
	sysInfo.BandwidthOut = roundToTwoDecimals(txBytesPerSec)

	return sysInfo, nil
}

// coverToGb 将字节数转换为GB, 保留两位小数
func coverToGb(b uint64) float64 {
	f := float64(b) / (1024 * 1024 * 1024)
	return roundToTwoDecimals(f)
}

// roundToTwoDecimals 将浮点数保留两位小数
func roundToTwoDecimals(f float64) float64 {
	return math.Round(f*100) / 100
}
