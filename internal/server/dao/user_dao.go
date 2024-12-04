package dao

import (
	"noah/internal/server/model"
	"time"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(i do.Injector) (UserDao, error) {
	d := UserDao{db: do.MustInvoke[*gorm.DB](i)}

	d.execInitSql()

	return d, nil
}

func (d UserDao) execInitSql() {
	u, err := d.QueryById(1)
	if err == nil && u.ID == 1 {
		return
	}
	sql := "insert into user(id, username, password, name, avatar) values (1, 'admin', '123456', '管理员', 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif')"
	d.db.Exec(sql)
}

func (d UserDao) QueryByUsername(username string) (user model.User, err error) {
	err = d.db.Where("username = ?", username).First(&user).Error
	return
}

func (d UserDao) QueryById(id uint) (user model.User, err error) {
	err = d.db.Where("id = ?", id).First(&user).Error
	return
}

func (d UserDao) Page(page, size int) (total int64, users []model.User, err error) {
	err = d.db.Model(&model.User{}).Count(&total).Error
	if err != nil {
		return
	}
	err = d.db.Scopes(Paginate(page, size)).Find(&users).Error
	return
}

func (d UserDao) UpdateToken(id uint, token string, refreshToken string, expireTime time.Time, refreshExpireTime time.Time) error {
	return d.db.Model(&model.User{}).Where("id = ?", id).Updates(model.User{
		Token:             token,
		RefreshToken:      refreshToken,
		ExpireTime:        expireTime,
		RefreshExpireTime: refreshExpireTime,
		LoginTime:         time.Now(),
	}).Error
}

func (d UserDao) UpdatePassword(id uint, password string) error {
	return d.db.Model(&model.User{}).Where("id = ?", id).Update("password", password).Error
}

func (d UserDao) QueryAll() (users []model.User, err error) {
	err = d.db.Find(&users).Error
	return
}
