package dao

import (
	"gorm.io/gorm"
	"time"
)

type ClientInfoDao struct {
	Db *gorm.DB
}

type ClientInfo struct {
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

func (ClientInfo) TableName() string {
	return "client_info"
}

func (d ClientInfoDao) Create(clientInfo ClientInfo) (err error) {
	err = d.Db.Create(&clientInfo).Error
	return err
}

func (d ClientInfoDao) GetByClientId(clientId uint, start time.Time, end time.Time) (clientInfoList []ClientInfo) {
	d.Db.Where("client_id = ?", clientId).Where("created_at BETWEEN ? AND ?", start, end).Find(&clientInfoList)
	return clientInfoList
}

func (d ClientInfoDao) Clean() {
	d.Db.Unscoped().Where("created_at < ?", time.Now().AddDate(0, 0, -1)).Delete(&ClientInfo{})
}
