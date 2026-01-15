package handler

import (
	"encoding/json"
	"fmt"
	"noah/pkg/packet"
	"runtime"
	"time"

	"noah/pkg/sys"

	"github.com/jinzhu/copier"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
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

func (h *InfoHandler) GetInfo() *packet.AgentInfo {
	info := &packet.AgentInfo{}

	basicInfo, _ := sys.GetBasicInfo()
	copier.Copy(info, basicInfo)

	hostInfo, _ := sys.GetHostInfo()
	copier.Copy(info, hostInfo)

	cpuInfo, _ := sys.GetCpuInfo()
	if cpuInfo != nil {
		copier.Copy(info, cpuInfo)
		js, _ := json.Marshal(cpuInfo.CpuInfo)
		if len(js) > 0 {
			info.CpuInfo = string(js)
		}
	}

	info.MemTotal = sys.GetTotalMemory()

	info.DiskTotal = sys.GetDiskTotal()

	return info
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
