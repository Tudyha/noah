package dao

import (
	"context"
	"noah/internal/model"
	"noah/pkg/enum"
	"noah/pkg/request"
	"time"

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

func (c *clientDao) UpdateStatus(ctx context.Context, clientID uint64, status enum.ClientStatus) error {
	return c.db.WithContext(ctx).Model(&model.Client{}).Where("id = ?", clientID).Update("status", status).Error
}

func (c *clientDao) Delete(ctx context.Context, clientID uint64) error {
	return c.db.WithContext(ctx).Delete(&model.Client{}, clientID).Error

}

func (c *clientDao) GetByID(ctx context.Context, clientID uint64) (*model.Client, error) {
	var client model.Client
	return &client, c.db.WithContext(ctx).First(&client, clientID).Error
}

func (c *clientDao) GetByIDs(ctx context.Context, clientIDs []uint64) ([]*model.Client, error) {
	var clients []*model.Client
	return clients, c.db.WithContext(ctx).Model(&model.Client{}).Where("id in ?", clientIDs).Find(&clients).Error
}

func (c *clientDao) SaveClientStat(ctx context.Context, stat *model.ClientStat) error {
	return c.db.WithContext(ctx).Model(&model.ClientStat{}).Save(stat).Error
}

func (c *clientDao) GetClientStat(ctx context.Context, clientID uint64, start time.Time, end time.Time) (any, error) {
	var stats []*model.ClientStat
	return stats, c.db.WithContext(ctx).Model(&model.ClientStat{}).Where("client_id = ?", clientID).Where("created_at >= ?", start).Where("created_at <= ?", end).Find(&stats).Error
}

func (c *clientDao) GetBySessionID(ctx context.Context, sessionID string) (*model.Client, error) {
	var model model.Client
	return &model, c.db.WithContext(ctx).Model(&model).Where("session_id = ?", sessionID).First(&model).Error
}
