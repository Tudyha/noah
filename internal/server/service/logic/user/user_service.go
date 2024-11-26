package user

import (
	"noah/internal/server/dao"
	"noah/internal/server/model"
	"noah/pkg/errcode"
	"time"

	"github.com/samber/do/v2"
)

type userService struct {
	userDao dao.UserDao
}

func NewUserService(i do.Injector) (userService, error) {
	return userService{
		userDao: do.MustInvoke[dao.UserDao](i),
	}, nil
}

func (s userService) Login(username string, password string) (user model.User, err error) {
	user, err = s.userDao.QueryByUsername(username)
	if err != nil {
		return user, errcode.ErrPasswordNotMatch
	}
	if user.Password != password {
		return user, errcode.ErrPasswordNotMatch
	}
	return user, nil
}

func (s userService) UpdateToken(id uint, token string, refreshToken string, expireTime time.Time, refreshExpireTime time.Time) error {
	return s.userDao.UpdateToken(id, token, refreshToken, expireTime, refreshExpireTime)
}

func (s userService) GetUser(userId uint) (user model.User, err error) {
	user, err = s.userDao.QueryById(userId)
	return user, err
}

func (s userService) GetUserPage(page, size int) (total int64, users []model.User, err error) {
	total, users, err = s.userDao.Page(page, size)
	return total, users, err
}

func (s userService) UpdatePassword(id uint, password string) error {
	return s.userDao.UpdatePassword(id, password)
}

func (s userService) GetUserList() (users []model.User, err error) {
	return s.userDao.QueryAll()
}
