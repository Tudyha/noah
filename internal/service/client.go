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

type clientService struct {
	clientDao dao.ClientDao
	spaceDao  dao.WorkSpaceDao
}

func newClientService() ClientService {
	return &clientService{
		clientDao: dao.GetClientDao(),
		spaceDao:  dao.GetWorkSpaceDao(),
	}
}

func (c *clientService) Connect(ctx context.Context, client *model.Client) error {
	old, err := c.clientDao.GetByDeviceID(ctx, client.DeviceID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	client.ID = old.ID
	client.Status = enum.ClientStatusOnline
	client.LastOnlineTime = time.Now()
	return c.clientDao.Create(ctx, client)
}

func (c *clientService) VerifySign(ctx context.Context, appID uint64, sign string) error {
	app, err := c.spaceDao.GetAppByAppID(ctx, appID)
	if err != nil {
		return err
	}
	if !utils.VerifySignature(app.ID, app.Secret, sign) {
		return errors.New("sign verify error")
	}
	return nil
}

func (c *clientService) GetPage(ctx context.Context, appID uint64, query request.ClientQueryRequest) (*response.Page[response.ClientResponse], error) {
	clients, total, err := c.clientDao.GetPage(ctx, appID, query)
	if err != nil {
		return nil, err
	}
	var list []response.ClientResponse
	copier.Copy(&list, clients)
	return &response.Page[response.ClientResponse]{
		Total: total,
		List:  list,
	}, nil
}

func (c *clientService) Disconnect(ctx context.Context, clientID uint64) error {
	return c.clientDao.UpdateStatus(ctx, clientID, enum.ClientStatusOffline)
}

func (c *clientService) Delete(ctx context.Context, clientID uint64) (*model.Client, error) {
	old, err := c.clientDao.GetByID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	return old, c.clientDao.Delete(ctx, clientID)
}

func (c *clientService) SaveClientStat(ctx context.Context, stat *model.ClientStat) error {
	return c.clientDao.SaveClientStat(ctx, stat)
}

func (c *clientService) GetClientStat(ctx context.Context, clientID uint64, start time.Time, end time.Time) ([]*response.ClientStatResponse, error) {
	stats, err := c.clientDao.GetClientStat(ctx, clientID, start, end)
	if err != nil {
		return nil, err
	}
	var list []*response.ClientStatResponse
	copier.Copy(&list, stats)
	return list, nil
}
