package dao

import (
	"context"
	"noah/internal/model"
	"noah/pkg/enum"
	"noah/pkg/request"
	"time"

	"gorm.io/gorm"
)

type agentDao struct {
	db *gorm.DB
}

func newAgentDao(db *gorm.DB) AgentDao {
	return &agentDao{
		db: db,
	}
}

func (c *agentDao) Create(ctx context.Context, agent *model.Agent) error {
	return c.db.WithContext(ctx).Save(agent).Error
}

func (c *agentDao) GetByDeviceID(ctx context.Context, deviceID string) (*model.Agent, error) {
	var agent model.Agent
	return &agent, c.db.WithContext(ctx).Where("device_id = ?", deviceID).First(&agent).Error
}

func (c *agentDao) GetPage(ctx context.Context, appID uint64, query request.AgentQueryRequest) ([]*model.Agent, int64, error) {
	var agents []*model.Agent
	var total int64
	db := c.db.WithContext(ctx).Model(&model.Agent{})
	if appID != 0 {
		db = db.Where("app_id = ?", appID)
	}
	if query.Status != 0 {
		db = db.Where("status = ?", query.Status)
	}
	if query.Hostname != "" {
		db = db.Where("hostname like ?", "%"+query.Hostname+"%")
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	db = db.Order("id desc")
	return agents, total, db.Scopes(Paginate(query.PageQuery)).Find(&agents).Error
}

func (c *agentDao) UpdateStatus(ctx context.Context, agentID uint64, status enum.AgentStatus) error {
	return c.db.WithContext(ctx).Model(&model.Agent{}).Where("id = ?", agentID).Update("status", status).Error
}

func (c *agentDao) Delete(ctx context.Context, agentID uint64) error {
	return c.db.WithContext(ctx).Delete(&model.Agent{}, agentID).Error

}

func (c *agentDao) GetByID(ctx context.Context, agentID uint64) (*model.Agent, error) {
	var agent model.Agent
	return &agent, c.db.WithContext(ctx).First(&agent, agentID).Error
}

func (c *agentDao) GetByIDs(ctx context.Context, agentIDs []uint64) ([]*model.Agent, error) {
	var agents []*model.Agent
	return agents, c.db.WithContext(ctx).Model(&model.Agent{}).Where("id in ?", agentIDs).Find(&agents).Error
}

func (c *agentDao) SaveAgentMetric(ctx context.Context, stat *model.AgentMetric) error {
	return c.db.WithContext(ctx).Model(&model.AgentMetric{}).Save(stat).Error
}

func (c *agentDao) GetAgentMetric(ctx context.Context, agentID uint64, start time.Time, end time.Time) (any, error) {
	var stats []*model.AgentMetric
	return stats, c.db.WithContext(ctx).Model(&model.AgentMetric{}).Where("agent_id = ?", agentID).Where("created_at >= ?", start).Where("created_at <= ?", end).Find(&stats).Error
}

func (c *agentDao) GetBySessionID(ctx context.Context, sessionID string) (*model.Agent, error) {
	var model model.Agent
	return &model, c.db.WithContext(ctx).Model(&model).Where("session_id = ?", sessionID).First(&model).Error
}

func (c *agentDao) CountByAppID(ctx context.Context, appID uint64) (online int64, offline int64, err error) {
	err = c.db.WithContext(ctx).Model(&model.Agent{}).Where("app_id = ? AND status = ?", appID, enum.AgentStatusOnline).Count(&online).Error
	if err != nil {
		return
	}
	err = c.db.WithContext(ctx).Model(&model.Agent{}).Where("app_id = ? AND status = ?", appID, enum.AgentStatusOffline).Count(&offline).Error
	return
}
