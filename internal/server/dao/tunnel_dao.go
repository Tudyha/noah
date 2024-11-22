package dao

import (
	"noah/internal/server/model"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type TunnelDao struct {
	db *gorm.DB
}

func NewTunnelDao(i do.Injector) (TunnelDao, error) {
	return TunnelDao{db: do.MustInvoke[*gorm.DB](i)}, nil
}

func (d TunnelDao) Save(tunnel model.Tunnel) (id uint, err error) {
	result := d.db.Create(&tunnel)
	if result.Error != nil {
		return 0, result.Error
	}
	return tunnel.ID, nil
}

func (d TunnelDao) GetById(tunnelId uint) (tunnel model.Tunnel, err error) {
	err = d.db.Where("id = ?", tunnelId).First(&tunnel).Error
	return tunnel, err
}

func (d TunnelDao) Delete(id uint) error {
	d.db.Unscoped().Delete(&model.Tunnel{}, id)
	return nil
}

func (d TunnelDao) List(clientId uint) (tunnels []model.Tunnel, err error) {
	if clientId == 0 {
		err = d.db.Find(&tunnels).Error
		return tunnels, err
	}
	err = d.db.Where("client_id = ?", clientId).Find(&tunnels).Error
	return tunnels, err
}

func (d TunnelDao) UpdateStatus(id uint, status uint8, failReason string) error {
	return d.db.Model(&model.Tunnel{}).Where("id = ?", id).Updates(model.Tunnel{Status: status, FailReason: failReason}).Error
}
