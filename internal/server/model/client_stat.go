package model

import "gorm.io/gorm"

type ClientStat struct {
	gorm.Model
	ClientID          uint    `gorm:"comment:客户端id;not null"`
	CpuUsage          float64 `gorm:"comment:CPU使用率"`
	MemoryTotal       uint64  `gorm:"comment:总内存"`
	MemoryUsed        uint64  `gorm:"comment:已用内存"`
	MemoryFree        uint64  `gorm:"comment:空闲内存"`
	MemoryUsedPercent float64 `gorm:"comment:内存使用百分比"`
	MemoryAvailable   uint64  `gorm:"comment:可用内存"`
	DiskTotal         uint64  `gorm:"comment:磁盘总量"`
	DiskFree          uint64  `gorm:"comment:空闲磁盘空间"`
	DiskUsed          uint64  `gorm:"comment:已用磁盘空间"`
	BandwidthIn       float64 `gorm:"comment:带宽入"`
	BandwidthOut      float64 `gorm:"comment:带宽出"`
}

func (ClientStat) TableName() string {
	return "client_stat"
}
