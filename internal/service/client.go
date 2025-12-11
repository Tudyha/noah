package service

import (
	"context"
	"errors"
	"noah/internal/dao"
	"noah/internal/model"
	"noah/pkg/enum"
	"noah/pkg/utils"
	"time"

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

func (c *clientService) Create(ctx context.Context, client *model.Client) error {
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
