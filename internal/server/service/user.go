package service

import (
	"noah/internal/server/model"
	"time"
)

type IUserService interface {
	Login(username string, password string) (model.User, error)
	GetUser(id uint) (model.User, error)
	GetUserPage(page, size int) (int64, []model.User, error)
	UpdateToken(id uint, token string, refreshToken string, expireTime time.Time, refreshExpireTime time.Time) error
	UpdatePassword(id uint, password string) error
}
