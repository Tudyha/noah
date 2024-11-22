package dao

import (
	"noah/internal/server/model"
	"time"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type ClientStatDao struct {
	Db *gorm.DB
}

func NewClientStatDao(i do.Injector) (ClientStatDao, error) {
	return ClientStatDao{
		Db: do.MustInvoke[*gorm.DB](i),
	}, nil
}

func (d ClientStatDao) Create(clientStat model.ClientStat) (err error) {
	err = d.Db.Create(&clientStat).Error
	return err
}

func (d ClientStatDao) GetByClientId(clientId uint, start time.Time, end time.Time) (clientStatList []model.ClientStat) {
	d.Db.Where("client_id = ?", clientId).Where("created_at BETWEEN ? AND ?", start, end).Find(&clientStatList)
	return
}

func (d ClientStatDao) Clean() {
	d.Db.Unscoped().Where("created_at < ?", time.Now().AddDate(0, 0, -1)).Delete(&model.ClientStat{})
}
