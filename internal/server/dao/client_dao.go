package dao

import (
	"fmt"
	"noah/internal/server/model"
	"noah/pkg/enum"
	"strings"
	"time"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type ClientDao struct {
	Db *gorm.DB
}

func NewClientDao(i do.Injector) (ClientDao, error) {
	return ClientDao{
		Db: do.MustInvoke[*gorm.DB](i),
	}, nil
}

func (d ClientDao) Save(client model.Client) (id uint, err error) {
	//状态默认在线
	client.Status = enum.DEVICE_ONLINE
	client.LastOnlineTime = time.Now()

	result := d.Db.Create(&client)
	if result.Error != nil {
		return 0, result.Error
	}
	return client.ID, nil
}

func (d ClientDao) Update(Client model.Client) (err error) {
	Client.Status = enum.DEVICE_ONLINE
	Client.LastOnlineTime = time.Now()
	result := d.Db.Updates(&Client)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d ClientDao) GetByMacAddress(macAddress string) (Client model.Client) {
	d.Db.Where("mac_address = ?", macAddress).First(&Client)
	return Client
}

func (d ClientDao) UpdateStatus(id uint, status int8) {
	d.Db.Model(&model.Client{}).Where("id = ?", id).Updates(model.Client{Status: status, LastOnlineTime: time.Now()})
}

func (d ClientDao) ScheduleUpdateStatus() {
	d.Db.Model(&model.Client{}).Where("last_online_time < ?", time.Now().Add(-10*time.Minute)).Updates(model.Client{Status: enum.DEVICE_OFFLINE})
}

func (d ClientDao) Page(hostname string, status int8, page, size int) (total int64, Clients []model.Client) {
	qw := d.Db

	// 处理 hostname 查询条件
	if hostname != "" {
		hostnameCond := fmt.Sprintf("hostname LIKE '%%%s%%'", strings.ReplaceAll(hostname, "'", "''"))
		qw = qw.Where(hostnameCond)
	}

	// 处理 status 查询条件
	if status != 0 {
		qw = qw.Where("status = ?", status)
	}

	// 分页查询
	err := qw.Scopes(Paginate(page, size)).Find(&Clients).Error
	if err != nil {
		// 处理数据库操作错误
		fmt.Println("Database error:", err)
		return 0, nil
	}

	// 计算总数
	qw.Model(&model.Client{}).Count(&total)

	return total, Clients
}

func (d ClientDao) GetById(id uint) (client model.Client, err error) {
	err = d.Db.Where("id = ?", id).First(&client).Error
	return client, err
}

func (d ClientDao) Delete(id uint) error {
	d.Db.Unscoped().Delete(&model.Client{}, id)
	return nil
}

func (d ClientDao) Count() (online int64, offline int64) {
	d.Db.Where("status = ?", enum.DEVICE_ONLINE).Model(&model.Client{}).Count(&online)
	d.Db.Where("status = ?", enum.DEVICE_OFFLINE).Model(&model.Client{}).Count(&offline)
	return online, offline
}
