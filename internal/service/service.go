package service

import (
	"context"
	"noah/internal/dao"
	"noah/internal/model"
	"noah/pkg/enum"
	"noah/pkg/request"
	"noah/pkg/response"
	"sync"
)

var (
	once sync.Once

	userServiceInstance UserService
	authServiceInstance AuthService
	smsServiceInstance  SmsService
)

// UserService 用户服务接口
type UserService interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint64) (*response.UserResponse, error)
}

// AuthService 认证服务接口
type AuthService interface {
	Login(ctx context.Context, req request.LoginRequest) (response.LoginResponse, error)
	Logout(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
	ValidateToken(ctx context.Context, token string) (uint64, error)
}

type SmsService interface {
	SendCode(ctx context.Context, smsCodeType enum.SmsCodeType, phone string) error
	VerifyCode(ctx context.Context, smsCodeType enum.SmsCodeType, phone, code string) error
}

type ClientService interface {
	Create(ctx context.Context) error
}

func Init() error {
	once.Do(func() {
		// 初始化服务
		smsServiceInstance = newSmsService()
		userServiceInstance = newUserService(dao.GetUserDao(), dao.GetWorkSpaceDao())
		authServiceInstance = newAuthService(dao.GetUserDao(), smsServiceInstance, userServiceInstance)
	})
	return nil
}

func GetUserService() UserService {
	return userServiceInstance
}

func GetAuthService() AuthService {
	return authServiceInstance
}
