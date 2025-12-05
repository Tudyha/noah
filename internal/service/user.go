package service

import (
	"context"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/samber/lo"

	"noah/internal/dao"
	"noah/internal/model"
	"noah/pkg/errcode"
	"noah/pkg/response"
)

type userService struct {
	userDao      dao.UserDao
	workSpaceDao dao.WorkSpaceDao
}

func newUserService(userDao dao.UserDao, workSpaceDao dao.WorkSpaceDao) UserService {
	return &userService{
		userDao:      userDao,
		workSpaceDao: workSpaceDao,
	}
}

func (s *userService) Create(ctx context.Context, user *model.User) error {
	if err := s.userDao.Create(ctx, user); err != nil {
		return err
	}

	user.Nickname = fmt.Sprintf("用户%06d", user.ID)
	if err := s.userDao.Update(ctx, user); err != nil {
		return err
	}

	// 创建默认空间
	space, err := s.workSpaceDao.Create(ctx, fmt.Sprintf("%s的空间", user.Nickname), "")
	if err != nil {
		return err
	}
	if err := s.workSpaceDao.CreateSpaceUser(ctx, space.ID, user.ID, 1); err != nil {
		return err
	}
	if err := s.workSpaceDao.CreateApp(ctx, space.ID, "默认应用", ""); err != nil {
		return err
	}

	return nil
}

func (s *userService) GetByID(ctx context.Context, id uint64) (*response.UserResponse, error) {
	user, err := s.userDao.FindByID(ctx, id)
	if err != nil {
		return nil, errcode.ErrUserNotFound
	}

	var userResp response.UserResponse
	copier.Copy(&userResp, user)
	userResp.WorkSpaceList, err = s.GetSpaceList(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &userResp, nil
}

func (s *userService) GetSpaceList(ctx context.Context, userID uint64) ([]*response.WorkSpaceResponse, error) {
	spaceList, err := s.workSpaceDao.GetByUserID(ctx, userID)
	if err != nil || len(spaceList) == 0 {
		return nil, err
	}

	spaceIds := lo.Map(spaceList, func(s *model.WorkSpace, index int) uint64 {
		return s.ID
	})
	appList, err := s.workSpaceDao.GetAppBySpaceIDs(ctx, spaceIds)
	if err != nil || len(appList) == 0 {
		return nil, err
	}

	appMap := lo.GroupByMap(appList, func(item *model.WorkSpaceApp) (uint64, *model.WorkSpaceApp) {
		return item.SpaceID, item
	})

	var list []*response.WorkSpaceResponse

	copier.Copy(&list, spaceList)

	for _, space := range list {
		apps := appMap[space.ID]
		if len(apps) > 0 {
			var appResp []*response.WorkSpaceAppResponse
			copier.Copy(&appResp, apps)
			space.AppList = appResp
		}
	}

	return list, nil
}
