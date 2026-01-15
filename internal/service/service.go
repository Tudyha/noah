package service

import (
	"context"
	"noah/internal/dao"
	"noah/internal/model"
	"noah/pkg/enum"
	"noah/pkg/request"
	"noah/pkg/response"
	"sync"
	"time"
)

var (
	once sync.Once

	userServiceInstance  UserService
	authServiceInstance  AuthService
	smsServiceInstance   SmsService
	agentServiceInstance AgentService
	workServiceInstance  WorkService
)

// UserService 用户服务接口
type UserService interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint64) (*response.UserResponse, error)
}

// AuthService 认证服务接口
type AuthService interface {
	Login(ctx context.Context, req request.LoginRequest) (response.LoginResponse, error)
	Register(ctx context.Context, req request.RegisterRequest) error
	SendCode(ctx context.Context, req request.SendCodeRequest) error
	Logout(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
	ValidateToken(ctx context.Context, token string) (uint64, error)
}

type SmsService interface {
	SendCode(ctx context.Context, smsCodeType enum.SmsCodeType, phone string) error
	VerifyCode(ctx context.Context, smsCodeType enum.SmsCodeType, phone, code string) error
}

type WorkService interface {
	GetAppByAppID(ctx context.Context, appID uint64) (*model.WorkSpaceApp, error)
}

type AgentService interface {
	VerifySign(ctx context.Context, appID uint64, sign string) error
	Connect(ctx context.Context, agent *model.Agent) error
	GetPage(ctx context.Context, appID uint64, query request.AgentQueryRequest) (*response.Page[response.AgentResponse], error)
	Disconnect(ctx context.Context, agentID uint64) error
	Delete(ctx context.Context, agentID uint64) (*model.Agent, error)
	SaveAgentMetric(ctx context.Context, sessionID string, agentMetric *model.AgentMetric) error
	GetAgentMetric(ctx context.Context, agentID uint64, start time.Time, end time.Time) ([]*response.AgentMetricResponse, error)
	GetByID(ctx context.Context, agentID uint64) (*model.Agent, error)
	GetByIDs(ctx context.Context, agentIDs []uint64) ([]*model.Agent, error)
	CountByAppID(ctx context.Context, appID uint64) (online int64, offline int64, err error)
}

func Init() error {
	once.Do(func() {
		// 初始化服务
		smsServiceInstance = newSmsService()
		userServiceInstance = newUserService(dao.GetUserDao(), dao.GetWorkSpaceDao())
		authServiceInstance = newAuthService(dao.GetUserDao(), smsServiceInstance, userServiceInstance)
		agentServiceInstance = newAgentService()
		workServiceInstance = newWorkService()
	})
	return nil
}

func GetUserService() UserService {
	return userServiceInstance
}

func GetAuthService() AuthService {
	return authServiceInstance
}

func GetAgentService() AgentService {
	return agentServiceInstance
}

func GetWorkService() WorkService {
	return workServiceInstance
}
