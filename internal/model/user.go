package model

import (
	"time"
)

type User struct {
	BaseModel

	Username          string    `gorm:"column:username;type:varchar(64);uniqueIndex;not null;comment:用户名"`
	Phone             string    `gorm:"column:phone;type:varchar(64);uniqueIndex;not null;comment:手机号"`
	Email             string    `gorm:"column:email;type:varchar(64);comment:邮箱"`
	Password          string    `gorm:"column:password;type:varchar(255);comment:密码"`
	Nickname          string    `gorm:"column:nickname;type:varchar(64);comment:昵称"`
	Avatar            string    `gorm:"column:avatar;type:varchar(255);comment:头像"`
	Token             string    `gorm:"column:token;type:varchar(255);comment:访问令牌"`
	RefreshToken      string    `gorm:"column:refresh_token;type:varchar(255);comment:刷新令牌"`
	ExpireTime        time.Time `gorm:"column:expire_time;type:datetime;comment:过期时间"`
	RefreshExpireTime time.Time `gorm:"column:refresh_expire_time;type:datetime;comment:刷新令牌过期时间"`
	LoginTime         time.Time `gorm:"column:login_time;type:datetime;comment:登录时间"`
	Status            int       `gorm:"column:status;type:tinyint;not null;default:1;comment:状态(0-禁用 1-启用)"`
}

func (User) TableName() string {
	return "user"
}
