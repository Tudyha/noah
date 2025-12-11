package model

import "time"

type Client struct {
	BaseModel

	AppId           uint64 `gorm:"column:app_id;not null"`
	DeviceID        string `gorm:"column:device_id;uniqueIndex"`
	OsType          int    `gorm:"column:os_type;type:int(11);not null"`
	Hostname        string `gorm:"column:hostname"`
	Username        string `gorm:"column:username"`
	Gid             string `gorm:"column:gid"`
	UID             string `gorm:"column:uid"`
	OsName          string `gorm:"column:os_name"`
	OsArch          string `gorm:"column:os_arch"`
	RemoteIP        string `gorm:"column:remote_ip;not null"`
	RemoteIpCountry string `gorm:"column:remote_ip_country"`
	LocalIP         string `gorm:"column:local_ip"`
	Port            string `gorm:"column:port"`
	Uptime          uint64 `gorm:"column:uptime"`
	BootTime        uint64 `gorm:"column:boot_time"`
	OS              string `gorm:"column:os"`
	Platform        string `gorm:"column:platform"`
	PlatformFamily  string `gorm:"column:platform_family"`
	PlatformVersion string `gorm:"column:platform_version"`
	KernelVersion   string `gorm:"column:kernel_version"`
	KernelArch      string `gorm:"column:kernel_arch"`
	HostID          string `gorm:"column:host_id"`
	CpuNum          int    `gorm:"column:cpu_num"`
	CpuInfo         string `gorm:"column:cpu_info;type:json"`
	MemTotal        uint64 `gorm:"column:mem_total"`
	DiskTotal       uint64 `gorm:"column:disk_total"`

	Status         int       `gorm:"column:status;type:int(11)"`
	LastOnlineTime time.Time `gorm:"column:last_online_time;type:datetime"`
}

// TableName table name
func (m *Client) TableName() string {
	return "client"
}
