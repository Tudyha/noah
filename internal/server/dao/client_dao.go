package dao

import (
	"fmt"
	"noah/internal/server/enum"
	"noah/internal/server/request"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ClientDao struct {
	Db *gorm.DB
}

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
	MemoryTotal    float64     `gorm:"comment:内存大小"`
	DiskTotal      float64     `gorm:"comment:磁盘大小"`
	RemoteIp       string      `gorm:"comment:公网ip"`
	LocalIp        string      `gorm:"comment:内网ip"`
	Port           string      `gorm:"comment:端口号"`
	Status         int8        `gorm:"default:0;comment:客户端状态 0-offline,1-online"`
	LastOnlineTime time.Time   `gorm:"comment:最后上线时间"`
}

func (Client) TableName() string {
	return "client"
}

func (d ClientDao) Save(client Client) (id uint, err error) {
	//状态默认在线
	client.Status = enum.DEVICE_ONLINE
	client.LastOnlineTime = time.Now()

	result := d.Db.Create(&client)
	if result.Error != nil {
		return 0, result.Error
	}
	return client.ID, nil
}

func (d ClientDao) Update(Client Client) (err error) {
	Client.Status = enum.DEVICE_ONLINE
	Client.LastOnlineTime = time.Now()
	result := d.Db.Updates(&Client)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d ClientDao) GetByMacAddress(macAddress string) (Client Client) {
	d.Db.Where("mac_address = ?", macAddress).First(&Client)
	return Client
}

func (d ClientDao) UpdateStatus(id uint, status int8) {
	d.Db.Model(&Client{}).Where("id = ?", id).Updates(Client{Status: status, LastOnlineTime: time.Now()})
}

func (d ClientDao) ScheduleUpdateStatus() {
	d.Db.Model(&Client{}).Where("last_online_time < ?", time.Now().Add(-10*time.Minute)).Updates(Client{Status: enum.DEVICE_OFFLINE})
}

func (d ClientDao) Page(query request.ListClientQueryReq) (total int64, Clients []Client) {
	qw := d.Db

	// 处理 hostname 查询条件
	if query.Hostname != "" {
		hostnameCond := fmt.Sprintf("hostname LIKE '%%%s%%'", strings.ReplaceAll(query.Hostname, "'", "''"))
		qw = qw.Where(hostnameCond)
	}

	// 处理 status 查询条件
	if query.Status != 0 {
		qw = qw.Where("status = ?", query.Status)
	}

	// 分页查询
	err := qw.Scopes(request.Paginate(&query.PageQuery)).Find(&Clients).Error
	if err != nil {
		// 处理数据库操作错误
		fmt.Println("Database error:", err)
		return 0, nil
	}

	// 计算总数
	qw.Model(&Client{}).Count(&total)

	return total, Clients
}

func (d ClientDao) GetById(id uint) (client Client, err error) {
	err = d.Db.Where("id = ?", id).First(&client).Error
	return client, err
}

func (d ClientDao) Delete(id uint) error {
	d.Db.Unscoped().Delete(&Client{}, id)
	return nil
}
