package information

import (
	"context"
	"fmt"
	"noah/pkg/mux/message"
	"noah/pkg/request"
	"noah/pkg/utils"
	"os"
	"os/user"
	"runtime"
	"time"

	"noah/pkg/utils/network"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/jinzhu/copier"
	"github.com/samber/do/v2"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

var (
	LastNetworkStatsTime time.Time
	LastNetworkBytesRecv uint64
	LastNetworkBytesSent uint64
)

type informationService struct {
}

func NewInformationService(i do.Injector) (informationService, error) {
	return informationService{}, nil
}

func (i informationService) LoadClientSpecs() (*request.CreateClientReq, error) {
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

	c, err := cpu.Info()
	if err != nil {
		return nil, err
	}
	cpu0 := c[0]

	memStats, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	usage, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}
	diskTotal := usage.Total

	return &request.CreateClientReq{
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
		MemoryTotal:  memStats.Total,
		DiskTotal:    diskTotal,
	}, nil
}

// GetSystemInfo 获取系统的基本信息，包括CPU、内存、磁盘和网络信息。
func (i informationService) GetSystemInfo() (*request.CreateClientStatReq, error) {
	var sysInfo = &request.CreateClientStatReq{}

	// 获取CPU信息
	_, err := cpu.Info()
	if err != nil {
		return nil, err
	}
	cpuUsage, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}
	sysInfo.CpuUsage = utils.RoundToTwoDecimals(cpuUsage[0])

	// 获取内存信息
	memStats, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	sysInfo.MemoryTotal = memStats.Total
	sysInfo.MemoryUsed = memStats.Used
	sysInfo.MemoryFree = memStats.Free
	sysInfo.MemoryAvailable = memStats.Available
	sysInfo.MemoryUsedPercent = utils.RoundToTwoDecimals(memStats.UsedPercent)

	// 获取磁盘信息
	usage, err := disk.Usage(".")
	if err != nil {
		return nil, err
	}
	sysInfo.DiskTotal = usage.Total
	sysInfo.DiskUsed = usage.Used
	sysInfo.DiskFree = usage.Free

	// 获取带宽信息
	bandwidth, err := net.IOCounters(false)
	if err != nil {
		return nil, fmt.Errorf("failed to get network bandwidth: %w", err)
	}
	var totalBytesRecv, totalBytesSent uint64
	for _, stats := range bandwidth {
		totalBytesRecv += stats.BytesRecv
		totalBytesSent += stats.BytesSent
	}
	if LastNetworkStatsTime.IsZero() {
		sysInfo.BandwidthIn = 0
		sysInfo.BandwidthOut = 0
	} else {
		sysInfo.BandwidthIn = float64(totalBytesRecv-LastNetworkBytesRecv) / time.Since(LastNetworkStatsTime).Seconds()
		sysInfo.BandwidthOut = float64(totalBytesSent-LastNetworkBytesSent) / time.Since(LastNetworkStatsTime).Seconds()
	}
	LastNetworkStatsTime = time.Now()
	LastNetworkBytesRecv = totalBytesRecv
	LastNetworkBytesSent = totalBytesSent

	return sysInfo, nil
}

func (i informationService) GetProcessList() ([]message.Process, error) {
	//获取系统进程
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}
	//
	var result []message.Process
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

		result = append(result, message.Process{
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

func (i informationService) KillProcess(pid int32) error {
	p, err := process.NewProcess(pid)
	if err != nil {
		return err
	}
	return p.Kill()
}

func (i informationService) GetNetworkInfo() (res []message.NetworkInfo, err error) {
	// 获取所有监听中的连接
	conns, err := net.Connections("-1")
	if err != nil {
		return nil, err
	}

	// 打印所有 TCP 和 TCP4 类型的连接
	copier.Copy(&res, &conns)
	return res, nil
}

func (i informationService) GetDockerContainerList() (res []message.DockerContainer, err error) {
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
