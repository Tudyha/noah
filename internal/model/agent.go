package model

import (
	"noah/pkg/enum"
	"time"
)

type Agent struct {
	BaseModel

	// 应用信息
	AppID          uint64           `gorm:"column:app_id;not null"`     // 应用ID
	Version        uint32           `gorm:"column:version"`             // agent版本号
	SessionID      string           `gorm:"column:session_id;not null"` // 连接id
	Status         enum.AgentStatus `gorm:"column:status"`              // 连接状态
	LastOnlineTime time.Time        `gorm:"column:last_online_time"`    // 最后一次连接时间

	// 设备信息
	DeviceID        string           `gorm:"column:device_id;uniqueIndex"`         // 设备ID
	OsType          enum.AgentOsType `gorm:"column:os_type;type:int(11);not null"` // 操作系统类型
	Hostname        string           `gorm:"column:hostname"`                      // 主机名
	Username        string           `gorm:"column:username"`                      // 用户名
	Gid             string           `gorm:"column:gid"`                           // 组ID
	UID             string           `gorm:"column:uid"`                           // 用户ID
	OsName          string           `gorm:"column:os_name"`                       // 操作系统名称
	OsArch          string           `gorm:"column:os_arch"`                       // 操作系统架构
	RemoteIP        string           `gorm:"column:remote_ip;not null"`            // 远程IP
	RemoteIpCountry string           `gorm:"column:remote_ip_country"`             // 远程IP所属国家
	LocalIP         string           `gorm:"column:local_ip"`                      // 本地IP
	Port            string           `gorm:"column:port"`                          // 端口
	Uptime          uint64           `gorm:"column:uptime"`                        // 运行时间
	BootTime        uint64           `gorm:"column:boot_time"`                     // 启动时间
	OS              string           `gorm:"column:os"`                            // 操作系统
	Platform        string           `gorm:"column:platform"`                      // 平台名称
	PlatformFamily  string           `gorm:"column:platform_family"`               // 平台
	PlatformVersion string           `gorm:"column:platform_version"`              // 平台版本
	KernelVersion   string           `gorm:"column:kernel_version"`                // 内核版本
	KernelArch      string           `gorm:"column:kernel_arch"`                   // 内核架构
	HostID          string           `gorm:"column:host_id"`                       // 主机ID
	CpuNum          int              `gorm:"column:cpu_num"`                       // cpu核数
	CpuInfo         string           `gorm:"column:cpu_info;type:json"`            // cpu信息
	MemTotal        uint64           `gorm:"column:mem_total"`                     // 内存总大小
	DiskTotal       uint64           `gorm:"column:disk_total"`                    // 磁盘总大小
}

func (m *Agent) TableName() string {
	return "agent"
}

type AgentMetric struct {
	BaseModel

	AgentID uint64 `gorm:"column:agent_id;not null"`

	// cpu数据
	CpuPercent float64 `gorm:"column:cpu_percent"`

	// 内存数据
	MemAvailable   uint64  `gorm:"column:mem_available"`
	MemUsed        uint64  `gorm:"column:mem_used"`
	MemUsedPercent float64 `gorm:"column:mem_used_percent"`
	MemFree        uint64  `gorm:"column:mem_free"`

	// 磁盘数据
	DiskUsage string `gorm:"column:disk_usage"`

	// 网络数据
	NetBytesSent float64 `gorm:"column:bytesSent"`
	NetBytesRecv float64 `gorm:"column:bytesRecv"`
}

func (m *AgentMetric) TableName() string {
	return "agent_metric"
}
