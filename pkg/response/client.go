package response

import (
	"noah/pkg/enum"
	"time"
)

type ClientResponse struct {
	ID              uint64            `json:"id"`
	AppID           uint64            `json:"app_id"`
	DeviceID        string            `json:"device_id"`
	OsType          enum.ClientOsType `json:"os_type"`
	Hostname        string            `json:"hostname"`
	Username        string            `json:"username"`
	Gid             string            `json:"gid"`
	UID             string            `json:"uid"`
	OsName          string            `json:"os_name"`
	OsArch          string            `json:"os_arch"`
	RemoteIP        string            `json:"remote_ip"`
	RemoteIpCountry string            `json:"remote_ip_country"`
	LocalIP         string            `json:"local_ip"`
	Port            string            `json:"port"`
	Uptime          uint64            `json:"uptime"`
	BootTime        uint64            `json:"boot_time"`
	OS              string            `json:"os"`
	Platform        string            `json:"platform"`
	PlatformFamily  string            `json:"platform_family"`
	PlatformVersion string            `json:"platform_version"`
	KernelVersion   string            `json:"kernel_version"`
	KernelArch      string            `json:"kernel_arch"`
	HostID          string            `json:"host_id"`
	CpuNum          int               `json:"cpu_num"`
	CpuInfo         string            `json:"cpu_info"`
	MemTotal        uint64            `json:"mem_total"`
	DiskTotal       uint64            `json:"disk_total"`
	ConnID          uint64            `json:"conn_id"`
	Status          enum.ClientStatus `json:"status"`
	LastOnlineTime  time.Time         `json:"last_online_time"`
	CreatedAt       *time.Time        `json:"created_at"`
	UpdatedAt       *time.Time        `json:"updated_at"`
}

type ClientBindResponse struct {
	MacBind     string `json:"mac_bind"`
	WindowsBind string `json:"windows_bind"`
	LinuxBind   string `json:"linux_bind"`
}

type ClientStatResponse struct {
	MemAvailable   uint64          `json:"mem_available"`
	MemUsed        uint64          `json:"mem_used"`
	MemUsedPercent float64         `json:"mem_used_percent"`
	MemFree        uint64          `json:"mem_free"`
	CpuPercent     float64         `json:"cpu_percent"`
	DiskUsage      []DiskUsageStat `json:"disk_usage"`
	NetBytesSent   float64         `json:"bytesSent"`
	NetBytesRecv   float64         `json:"bytesRecv"`
	CreatedAt      *time.Time      `json:"created_at"`
}

type DiskUsageStat struct {
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
