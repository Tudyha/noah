package information

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/jinzhu/copier"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"math"
	"noah/client/app/entitie"
	"noah/client/app/service"
	"os"
	"os/user"
	"runtime"

	"github.com/docker/docker/client"
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
	user, err := user.Current()
	if err != nil {
		return nil, err
	}
	macAddress, err := network.GetMacAddress()
	if err != nil {
		return nil, err
	}
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}
	cpu0 := cpuInfo[0]

	memStats, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	usage, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}
	diskTotal := coverToGb(usage.Total)

	return &entitie.Client{
		Hostname:     hostname,
		Username:     user.Username,
		Gid:          user.Gid,
		Uid:          user.Uid,
		OSName:       runtime.GOOS,
		OSArch:       runtime.GOARCH,
		MacAddress:   macAddress,
		IPAddress:    network.GetLocalIP(),
		CpuCores:     cpu0.Cores,
		CpuModelName: cpu0.ModelName,
		CpuFamily:    cpu0.Family,
		MemoryTotal:  coverToGb(memStats.Total),
		DiskTotal:    diskTotal,
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

func (i Service) GetProcessList() ([]entitie.Process, error) {
	//获取系统进程
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}
	//
	var result []entitie.Process
	for _, proc := range processes {
		if proc.Pid == 0 {
			fmt.Println("pid = 0 进程不显示")
			continue
		}
		name, _ := proc.Name()
		uid, _ := proc.Uids()
		gid, _ := proc.Gids()
		cmdline, _ := proc.Cmdline()
		username, _ := proc.Username()
		cpu, _ := proc.Percent(0)
		m, _ := proc.MemoryInfo()
		m_rss := uint64(0)
		if m != nil {
			m_rss = m.RSS
		}
		createTime, _ := proc.CreateTime()

		result = append(result, entitie.Process{
			Pid:        proc.Pid,
			Name:       name,
			Username:   username,
			Uids:       uid,
			Gids:       gid,
			Command:    cmdline,
			Cpu:        cpu,
			Memory:     m_rss,
			CreateTime: createTime,
		})
	}
	return result, nil
}

func (i Service) KillProcess(pid int32) error {
	p, err := process.NewProcess(pid)
	if err != nil {
		return err
	}
	return p.Kill()
}

func (i Service) GetNetworkInfo() (res []entitie.NetworkInfo, err error) {
	// 获取所有监听中的连接
	conns, err := net.Connections("-1")
	if err != nil {
		return nil, err
	}

	// 打印所有 TCP 和 TCP4 类型的连接
	copier.Copy(&res, &conns)
	return res, nil
}

func (i Service) GetDockerContainerList() (res []entitie.DockerContainer, err error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("Error creating Docker client:", err)
		return nil, err
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		fmt.Println("Error listing Docker containers:", err)
		return
	}

	copier.Copy(&res, &containers)
	return res, nil
}
