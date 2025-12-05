package dao

import (
	"context"
	"noah/internal/database"
	"noah/internal/model"
	"sync"
)

var (
	once sync.Once

	userDaoInstance      UserDao
	workSpaceDaoInstance WorkSpaceDao
)

func Init() error {
	db := database.GetDB()
	once.Do(func() {
		userDaoInstance = newUserDao(db)
		workSpaceDaoInstance = newWorkSpaceDao(db)
	})
	return nil
}

func GetUserDao() UserDao {
	return userDaoInstance
}

func GetWorkSpaceDao() WorkSpaceDao {
	return workSpaceDaoInstance
}

type UserDao interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id uint64) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
}

type WorkSpaceDao interface {
	Create(ctx context.Context, name, description string) (*model.WorkSpace, error)
	CreateSpaceUser(ctx context.Context, spaceId, userId uint64, role int) error
	CreateApp(ctx context.Context, spaceId uint64, name, description string) error
	GetByUserID(ctx context.Context, userID uint64) ([]*model.WorkSpace, error)
	GetAppBySpaceIDs(ctx context.Context, spaceIDs []uint64) ([]*model.WorkSpaceApp, error)
}
