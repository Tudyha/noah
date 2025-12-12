package dao

import (
	"context"
	"noah/internal/model"
	"noah/pkg/request"

	"gorm.io/gorm"
)

type clientDao struct {
	db *gorm.DB
}

func newClientDao(db *gorm.DB) ClientDao {
	return &clientDao{
		db: db,
	}
}

func (c *clientDao) Create(ctx context.Context, client *model.Client) error {
	return c.db.WithContext(ctx).Save(client).Error
}

func (c *clientDao) GetByDeviceID(ctx context.Context, deviceID string) (*model.Client, error) {
	var client model.Client
	return &client, c.db.WithContext(ctx).Where("device_id = ?", deviceID).First(&client).Error
}

func (c *clientDao) GetPage(ctx context.Context, appID uint64, query request.ClientQueryRequest) ([]*model.Client, int64, error) {
	var clients []*model.Client
	var total int64
	db := c.db.WithContext(ctx).Model(&model.Client{})
	if appID != 0 {
		db = db.Where("app_id = ?", appID)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	db = db.Order("id desc")
	return clients, total, db.Scopes(Paginate(query.PageQuery)).Find(&clients).Error
}
