package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"noah/pkg/packet"
	"os"
	"os/user"
	"runtime"
	"time"

	"github.com/jinzhu/copier"
	"github.com/samber/lo"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
)

var (
	LastNetStatTime  = time.Now()
	LastNetBytesRecv uint64
	LastNetBytesSent uint64
)

type InfoHandler struct {
}

func (h *InfoHandler) GetInfo() *packet.ClientInfo {
	info := &packet.ClientInfo{}
	err := h.getBasicInfo(info)
	if err != nil {
		log.Println("获取基本信息失败:", err)
		return info
	}
	err = h.getHostInfo(info)
	if err != nil {
		log.Println("获取主机信息失败:", err)
		return info
	}

	err = h.getCPUInfo(info)
	if err != nil {
		log.Println("获取CPU信息失败:", err)
		return info
	}

	err = h.getMemoryInfo(info)
	if err != nil {
		log.Println("获取内存信息失败:", err)
		return info
	}

	err = h.getDiskInfo(info)
	if err != nil {
		log.Println("获取磁盘信息失败:", err)
		return info
	}

	return info
}

func (h *InfoHandler) getBasicInfo(info *packet.ClientInfo) (err error) {
	// 获取主机名
	hostname, err := os.Hostname()
	if err != nil {
		return
	}
	info.Hostname = hostname

	// 获取当前用户
	currentUser, err := user.Current()
	if err != nil {
		return
	}
	info.Username = currentUser.Username
	info.Uid = currentUser.Uid
	info.Gid = currentUser.Gid

	// 操作系统和架构
	info.Os = runtime.GOOS
	info.OsArch = runtime.GOARCH
	info.KernelArch = runtime.GOARCH

	return
}

func (h *InfoHandler) getHostInfo(info *packet.ClientInfo) (err error) {
	hostInfo, err := host.Info()
	if err != nil {
		return
	}
	// 填充主机信息
	info.OsName = hostInfo.OS
	info.Platform = hostInfo.Platform
	info.PlatformFamily = hostInfo.PlatformFamily
	info.PlatformVersion = hostInfo.PlatformVersion
	info.KernelVersion = hostInfo.KernelVersion
	info.HostId = hostInfo.HostID

	// 启动时间相关
	info.BootTime = uint64(hostInfo.BootTime)
	if hostInfo.Uptime > 0 {
		info.Uptime = uint64(hostInfo.Uptime)
	} else {
		info.Uptime = uint64(time.Now().Unix()) - hostInfo.BootTime
	}

	return
}

func (h *InfoHandler) getCPUInfo(info *packet.ClientInfo) (err error) {
	// CPU数量
	info.CpuNum = int32(runtime.NumCPU())

	// 详细CPU信息
	cpuInfos, err := cpu.Info()
	if err != nil {
		return
	}

	if len(cpuInfos) > 0 {
		cpuInfo := cpuInfos[0]

		js, _ := json.Marshal(cpuInfo)
		if len(js) > 0 {
			info.CpuInfo = string(js)
		}
	}
	return
}

func (h *InfoHandler) getMemoryInfo(info *packet.ClientInfo) (err error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return
	}

	info.MemTotal = memInfo.Total
	return
}

func (h *InfoHandler) getDiskInfo(info *packet.ClientInfo) (err error) {
	// 获取所有分区
	partitions, err := getDiskUsage()
	if err != nil {
		return
	}

	info.DiskTotal = lo.Reduce(partitions, func(v uint64, item *disk.UsageStat, _ int) uint64 {
		return v + item.Total
	}, 0)
	return
}

func getDiskUsage() ([]*disk.UsageStat, error) {
	switch runtime.GOOS {
	case "linux", "darwin":
		u, err := disk.Usage("/")
		return []*disk.UsageStat{u}, err
	case "windows":
		parts, err := disk.Partitions(false)
		if err != nil {
			return nil, err
		}
		var us []*disk.UsageStat
		for _, p := range parts {
			u, err := disk.Usage(p.Mountpoint)
			if err != nil {
				return nil, err
			}
			us = append(us, u)
		}
		return us, nil
	default:
		return nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func (h *InfoHandler) GetSystemStat() *packet.Ping {
	info := &packet.Ping{}
	v, _ := mem.VirtualMemory()
	if v != nil {
		info.MemFree = v.Free
		info.MemUsed = v.Used
		info.MemUsedPercent = v.UsedPercent
		info.MemAvailable = v.Available
	}

	cpuPercent, _ := cpu.Percent(0, false)
	if len(cpuPercent) > 0 {
		info.CpuPercent = cpuPercent[0]
	}

	du, _ := getDiskUsage()
	copier.Copy(&info.DiskUsage, du)

	ioCounters, _ := net.IOCounters(false)

	if len(ioCounters) > 0 {
		info.NetBytesSent = float64(ioCounters[0].BytesSent-LastNetBytesSent) / time.Since(LastNetStatTime).Seconds()
		info.NetBytesRecv = float64(ioCounters[0].BytesRecv-LastNetBytesRecv) / time.Since(LastNetStatTime).Seconds()
		LastNetBytesSent = ioCounters[0].BytesSent
		LastNetBytesRecv = ioCounters[0].BytesRecv
	}
	LastNetStatTime = time.Now()

	return info
}
