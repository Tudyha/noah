package service

import (
	"context"
	"noah/internal/dao"
	"noah/internal/model"
)

type workService struct {
	spaceDao dao.WorkSpaceDao
}

func newWorkService() WorkService {
	return &workService{
		spaceDao: dao.GetWorkSpaceDao(),
	}
}

func (w *workService) GetAppByAppID(ctx context.Context, appID uint64) (*model.WorkSpaceApp, error) {
	return w.spaceDao.GetAppByAppID(ctx, appID)
}
