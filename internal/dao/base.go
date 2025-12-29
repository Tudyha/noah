package dao

import (
	"context"
	"noah/internal/database"
	"noah/internal/model"
	"noah/pkg/enum"
	"noah/pkg/request"
	"sync"
	"time"

	"gorm.io/gorm"
)

var (
	once sync.Once

	userDaoInstance      UserDao
	workSpaceDaoInstance WorkSpaceDao
	clientDaoInstance    ClientDao
)

func Init() error {
	db := database.GetDB()
	once.Do(func() {
		userDaoInstance = newUserDao(db)
		workSpaceDaoInstance = newWorkSpaceDao(db)
		clientDaoInstance = newClientDao(db)
	})
	return nil
}

func GetUserDao() UserDao {
	return userDaoInstance
}

func GetWorkSpaceDao() WorkSpaceDao {
	return workSpaceDaoInstance
}

func GetClientDao() ClientDao {
	return clientDaoInstance
}

type UserDao interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id uint64) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByPhone(ctx context.Context, phone string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
}

type WorkSpaceDao interface {
	Create(ctx context.Context, name, description string) (*model.WorkSpace, error)
	CreateSpaceUser(ctx context.Context, spaceId, userId uint64, role int) error
	CreateApp(ctx context.Context, spaceId uint64, secret, name, description string) error
	GetByUserID(ctx context.Context, userID uint64) ([]*model.WorkSpace, error)
	GetAppBySpaceIDs(ctx context.Context, spaceIDs []uint64) ([]*model.WorkSpaceApp, error)
	GetAppByAppID(ctx context.Context, appID uint64) (*model.WorkSpaceApp, error)
}

type ClientDao interface {
	Create(ctx context.Context, client *model.Client) error
	GetByDeviceID(ctx context.Context, deviceID string) (*model.Client, error)
	GetPage(ctx context.Context, appID uint64, query request.ClientQueryRequest) ([]*model.Client, int64, error)
	UpdateStatus(ctx context.Context, clientID uint64, status enum.ClientStatus) error
	Delete(ctx context.Context, clientID uint64) error
	GetByID(ctx context.Context, clientID uint64) (*model.Client, error)
	SaveClientStat(ctx context.Context, stat *model.ClientStat) error
	GetClientStat(ctx context.Context, clientID uint64, start time.Time, end time.Time) (any, error)
	GetBySessionID(ctx context.Context, sessionID string) (*model.Client, error)
	GetByIDs(ctx context.Context, clientIDs []uint64) ([]*model.Client, error)
}

func Paginate(pageQuery request.PageQuery) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := pageQuery.Page
		if page <= 0 {
			page = 1
		}

		pageSize := pageQuery.Limit
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
