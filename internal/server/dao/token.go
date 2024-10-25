package dao

import (
	"time"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type TokenDao struct {
	Db *gorm.DB
}

func NewTokenDao(i do.Injector) (*TokenDao, error) {
	return &TokenDao{
		Db: do.MustInvoke[*gorm.DB](i),
	}, nil
}

type Token struct {
	AuthDomain        uint8     `gorm:"comment:授权域"`
	Token             string    `gorm:"primarykey"`
	RefreshToken      string    `gorm:"comment:refreshToken"`
	ExpireTime        time.Time `gorm:"comment:token过期时间"`
	RefreshExpireTime time.Time `gorm:"comment:refreshToken过期时间"`
}

func (Token) TableName() string {
	return "token"
}

func (d TokenDao) Save(token Token) (err error) {
	result := d.Db.Create(&token)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d TokenDao) Delete(token string) (err error) {
	result := d.Db.Where("token = ?", token).Delete(&Token{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d TokenDao) List() (tokenInfo []Token, err error) {
	result := d.Db.Find(&tokenInfo)
	if result.Error != nil {
		return tokenInfo, result.Error
	}
	return tokenInfo, nil
}
