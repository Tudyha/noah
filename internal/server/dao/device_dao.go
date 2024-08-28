package dao

import (
	"noah/internal/server/dto"
	"noah/internal/server/enum"
	"time"

	"gorm.io/gorm"
)

type DeviceDao struct {
	Db *gorm.DB
}

type Device struct {
	gorm.Model
	Hostname       string      `gorm:"comment:主机名"`
	Username       string      `gorm:"comment:用户名"`
	UserID         string      `gorm:"comment:用户id"`
	OsType         enum.OSType `gorm:"comment:操作系统类型"`
	OSName         string      `gorm:"comment:系统名称"`
	OSArch         string      `gorm:"comment:系统发行版本"`
	MacAddress     string      `gorm:"unique;comment:mac地址"`
	IPAddress      string      `gorm:"comment:ip地址"`
	Port           string      `gorm:"comment:端口号"`
	Status         int8        `gorm:"default:0;comment:设备状态 0-offline,1-online"`
	LastOnlineTime time.Time   `gorm:"comment:最后上线时间"`
}

func (Device) TableName() string {
	return "device"
}

func (d *DeviceDao) Save(device Device) (id uint, err error) {
	device.Status = enum.DEVICE_ONLINE
	device.LastOnlineTime = time.Now()
	result := d.Db.Create(&device)
	if result.Error != nil {
		return 0, result.Error
	}
	return device.ID, nil
}

func (d *DeviceDao) Update(device Device) (err error) {
	device.Status = enum.DEVICE_ONLINE
	device.LastOnlineTime = time.Now()
	result := d.Db.Updates(&device)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *DeviceDao) GetByMacAddress(macAddress string) (device Device) {
	d.Db.Where("mac_address = ?", macAddress).First(&device)
	return device
}

func (d *DeviceDao) UpdateStatus(id uint, status int8) {
	d.Db.Model(&Device{}).Where("id = ?", id).Updates(Device{Status: status, LastOnlineTime: time.Now()})
}

func (d *DeviceDao) ScheduleUpdateStatus() {
	d.Db.Model(&Device{}).Where("last_online_time < ?", time.Now().Add(-10*time.Minute)).Updates(Device{Status: enum.DEVICE_OFFLINE})
}

func (d *DeviceDao) Page(query dto.DeviceListQueryDto) (total int64, devices []Device) {
	qw := d.Db.Where("1 = 1")
	if query.Hostname != "" {
		qw.Where("hostname LIKE ?", "%"+query.Hostname+"%")
	}
	if query.Status != 0 {
		qw.Where("status = ?", query.Status)
	}

	qw.Scopes(dto.Paginate(&query.PageQuery)).Find(&devices).Count(&total)
	return total, devices
}

func (d *DeviceDao) GetById(id uint) (device Device) {
	d.Db.Where("id = ?", id).First(&device)
	return device
}

func (d *DeviceDao) Delete(id uint) error {
	d.Db.Unscoped().Delete(&Device{}, id)
	return nil
}
