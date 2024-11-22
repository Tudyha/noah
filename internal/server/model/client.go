package model

import (
	"noah/pkg/enum"
	"time"

	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	Hostname       string      `gorm:"comment:主机名"`
	Username       string      `gorm:"comment:用户名"`
	Gid            string      `gorm:"comment:组id"`
	Uid            string      `gorm:"comment:用户id"`
	OsType         enum.OSType `gorm:"comment:操作系统类型"`
	OSName         string      `gorm:"comment:系统名称"`
	OSArch         string      `gorm:"comment:系统发行版本"`
	MacAddress     string      `gorm:"unique;comment:mac地址"`
	CpuCores       int32       `gorm:"comment:cpu核心数"`
	CpuModelName   string      `gorm:"comment:cpuModelName"`
	CpuFamily      string      `gorm:"comment:CpuFamily"`
	MemoryTotal    uint64      `gorm:"comment:内存大小"`
	DiskTotal      uint64      `gorm:"comment:磁盘大小"`
	RemoteIp       string      `gorm:"comment:公网ip"`
	LocalIp        string      `gorm:"comment:内网ip"`
	Port           string      `gorm:"comment:端口号"`
	Status         int8        `gorm:"default:0;comment:客户端状态 0-offline,1-online"`
	LastOnlineTime time.Time   `gorm:"comment:最后上线时间"`
}

func (Client) TableName() string {
	return "client"
}
