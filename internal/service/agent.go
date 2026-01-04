package service

import (
	"context"
	"errors"
	"noah/internal/dao"
	"noah/internal/model"
	"noah/pkg/enum"
	"noah/pkg/request"
	"noah/pkg/response"
	"noah/pkg/utils"
	"time"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type agentService struct {
	agentDao dao.AgentDao
	spaceDao dao.WorkSpaceDao
}

func newAgentService() AgentService {
	return &agentService{
		agentDao: dao.GetAgentDao(),
		spaceDao: dao.GetWorkSpaceDao(),
	}
}

func (c *agentService) Connect(ctx context.Context, agent *model.Agent) error {
	old, err := c.agentDao.GetByDeviceID(ctx, agent.DeviceID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	agent.ID = old.ID
	agent.Status = enum.AgentStatusOnline
	agent.LastOnlineTime = time.Now()
	return c.agentDao.Create(ctx, agent)
}

func (c *agentService) VerifySign(ctx context.Context, appID uint64, sign string) error {
	app, err := c.spaceDao.GetAppByAppID(ctx, appID)
	if err != nil {
		return err
	}
	if !utils.VerifySignature(app.ID, app.Secret, sign) {
		return errors.New("sign verify error")
	}
	return nil
}

func (c *agentService) GetPage(ctx context.Context, appID uint64, query request.AgentQueryRequest) (*response.Page[response.AgentResponse], error) {
	agents, total, err := c.agentDao.GetPage(ctx, appID, query)
	if err != nil {
		return nil, err
	}
	var list []response.AgentResponse
	copier.Copy(&list, agents)
	return &response.Page[response.AgentResponse]{
		Total: total,
		List:  list,
	}, nil
}

func (c *agentService) Disconnect(ctx context.Context, agentID uint64) error {
	return c.agentDao.UpdateStatus(ctx, agentID, enum.AgentStatusOffline)
}

func (c *agentService) Delete(ctx context.Context, agentID uint64) (*model.Agent, error) {
	old, err := c.agentDao.GetByID(ctx, agentID)
	if err != nil {
		return nil, err
	}

	return old, c.agentDao.Delete(ctx, agentID)
}

func (c *agentService) SaveAgentMetric(ctx context.Context, sessionID string, agentMetric *model.AgentMetric) error {
	agent, err := c.agentDao.GetBySessionID(ctx, sessionID)
	if err != nil {
		return err
	}
	agentMetric.AgentID = agent.ID
	return c.agentDao.SaveAgentMetric(ctx, agentMetric)
}

func (c *agentService) GetAgentMetric(ctx context.Context, agentID uint64, start time.Time, end time.Time) ([]*response.AgentMetricResponse, error) {
	metrics, err := c.agentDao.GetAgentMetric(ctx, agentID, start, end)
	if err != nil {
		return nil, err
	}
	var list []*response.AgentMetricResponse
	copier.Copy(&list, metrics)
	return list, nil
}

func (c *agentService) GetByID(ctx context.Context, agentID uint64) (*model.Agent, error) {
	return c.agentDao.GetByID(ctx, agentID)
}

func (c *agentService) GetByIDs(ctx context.Context, agentIDs []uint64) ([]*model.Agent, error) {
	return c.agentDao.GetByIDs(ctx, agentIDs)
}
