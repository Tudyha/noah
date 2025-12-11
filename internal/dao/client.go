package dao

import (
	"context"
	"noah/internal/model"

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
