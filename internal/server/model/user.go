package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username          string    `gorm:"comment:用户名"`
	Password          string    `gorm:"comment:密码"`
	Name              string    `gorm:"comment:name"`
	Introduction      string    `gorm:"comment:简介"`
	Avatar            string    `gorm:"comment:头像"`
	Token             string    `gorm:"comment:token"`
	RefreshToken      string    `gorm:"comment:refreshToken"`
	ExpireTime        time.Time `gorm:"comment:expireTime"`
	RefreshExpireTime time.Time `gorm:"comment:refreshExpireTime"`
	LoginTime         time.Time `gorm:"comment:登录时间"`
}

func (User) TableName() string {
	return "user"
}
