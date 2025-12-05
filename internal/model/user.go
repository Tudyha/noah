package model

import (
	"time"
)

type User struct {
	BaseModel

	Username          string    `gorm:"column:username;type:varchar(64);uniqueIndex;not null;comment:用户名" json:"username"`
	Phone             string    `gorm:"column:phone;type:varchar(64);uniqueIndex;not null;comment:手机号" json:"phone"`
	Password          string    `gorm:"column:password;type:varchar(255);comment:密码" json:"password"`
	Nickname          string    `gorm:"column:nickname;type:varchar(64);comment:昵称" json:"nickname"`
	Avatar            string    `gorm:"column:avatar;type:varchar(255);comment:头像" json:"avatar"`
	Token             string    `gorm:"column:token;type:varchar(255);comment:访问令牌" json:"token"`
	RefreshToken      string    `gorm:"column:refresh_token;type:varchar(255);comment:刷新令牌" json:"refreshToken"`
	ExpireTime        time.Time `gorm:"column:expire_time;type:datetime;comment:过期时间" json:"expireTime"`
	RefreshExpireTime time.Time `gorm:"column:refresh_expire_time;type:datetime;comment:刷新令牌过期时间" json:"refreshExpireTime"`
	LoginTime         time.Time `gorm:"column:login_time;type:datetime;comment:登录时间" json:"loginTime"`
	Status            int       `gorm:"column:status;type:tinyint;not null;default:1;comment:状态(0-禁用 1-启用)" json:"status"`
}

func (User) TableName() string {
	return "user"
}
