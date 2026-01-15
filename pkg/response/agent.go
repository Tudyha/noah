package response

import (
	"noah/pkg/enum"
	"time"
)

type AgentResponse struct {
	ID             uint64           `json:"id"`
	AppID          uint64           `json:"app_id"`
	Status         enum.AgentStatus `json:"status"`
	LastOnlineTime time.Time        `json:"last_online_time"`
	CreatedAt      *time.Time       `json:"created_at"`
	UpdatedAt      *time.Time       `json:"updated_at"`
	Version        uint32           `json:"version"`
	VersionName    string           `json:"version_name"`

	DeviceID        string           `json:"device_id"`
	OsType          enum.AgentOsType `json:"os_type"`
	Hostname        string           `json:"hostname"`
	Username        string           `json:"username"`
	Gid             string           `json:"gid"`
	UID             string           `json:"uid"`
	OsName          string           `json:"os_name"`
	OsArch          string           `json:"os_arch"`
	RemoteIP        string           `json:"remote_ip"`
	RemoteIpCountry string           `json:"remote_ip_country"`
	LocalIP         string           `json:"local_ip"`
	Port            string           `json:"port"`
	Uptime          uint64           `json:"uptime"`
	BootTime        uint64           `json:"boot_time"`
	OS              string           `json:"os"`
	Platform        string           `json:"platform"`
	PlatformFamily  string           `json:"platform_family"`
	PlatformVersion string           `json:"platform_version"`
	KernelVersion   string           `json:"kernel_version"`
	KernelArch      string           `json:"kernel_arch"`
	HostID          string           `json:"host_id"`
	CpuNum          int              `json:"cpu_num"`
	CpuInfo         string           `json:"cpu_info"`
	MemTotal        uint64           `json:"mem_total"`
	DiskTotal       uint64           `json:"disk_total"`
}

type AgentBindResponse struct {
	MacBind     string `json:"mac_bind"`
	WindowsBind string `json:"windows_bind"`
	LinuxBind   string `json:"linux_bind"`
}

type AgentMetricResponse struct {
	CreatedAt *time.Time `json:"created_at"`

	// cpu数据
	CpuPercent float64 `json:"cpu_percent"`

	// 内存数据
	MemAvailable   uint64  `json:"mem_available"`
	MemUsed        uint64  `json:"mem_used"`
	MemUsedPercent float64 `json:"mem_used_percent"`
	MemFree        uint64  `json:"mem_free"`

	// 磁盘数据
	DiskUsage []DiskUsage `json:"disk_usage"`

	// 网络数据
	NetBytesSent float64 `json:"bytesSent"`
	NetBytesRecv float64 `json:"bytesRecv"`
}

type DiskUsage struct {
	Path              string  `json:"path"`
	Fstype            string  `json:"fstype"`
	Total             uint64  `json:"total"`
	Free              uint64  `json:"free"`
	Used              uint64  `json:"used"`
	UsedPercent       float64 `json:"used_percent"`
	InodesTotal       uint64  `json:"inodes_total"`
	InodesFree        uint64  `json:"inodes_free"`
	InodesUsed        uint64  `json:"inodes_used"`
	InodesUsedPercent float64 `json:"inodes_used_percent"`
}
