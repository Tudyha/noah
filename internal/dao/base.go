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
	agentDaoInstance     AgentDao
)

func Init() error {
	db := database.GetDB()
	once.Do(func() {
		userDaoInstance = newUserDao(db)
		workSpaceDaoInstance = newWorkSpaceDao(db)
		agentDaoInstance = newAgentDao(db)
	})
	return nil
}

func GetUserDao() UserDao {
	return userDaoInstance
}

func GetWorkSpaceDao() WorkSpaceDao {
	return workSpaceDaoInstance
}

func GetAgentDao() AgentDao {
	return agentDaoInstance
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

type AgentDao interface {
	Create(ctx context.Context, Agent *model.Agent) error
	GetByDeviceID(ctx context.Context, deviceID string) (*model.Agent, error)
	GetPage(ctx context.Context, appID uint64, query request.AgentQueryRequest) ([]*model.Agent, int64, error)
	UpdateStatus(ctx context.Context, AgentID uint64, status enum.AgentStatus) error
	Delete(ctx context.Context, AgentID uint64) error
	GetByID(ctx context.Context, AgentID uint64) (*model.Agent, error)
	SaveAgentMetric(ctx context.Context, agentMetric *model.AgentMetric) error
	GetAgentMetric(ctx context.Context, AgentID uint64, start time.Time, end time.Time) (any, error)
	GetBySessionID(ctx context.Context, sessionID string) (*model.Agent, error)
	GetByIDs(ctx context.Context, AgentIDs []uint64) ([]*model.Agent, error)
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
