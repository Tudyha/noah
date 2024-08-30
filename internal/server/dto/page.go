package dto

import "gorm.io/gorm"

type PageQuery struct {
	Page int `form:"page" default:"1"`
	Size int `form:"size" default:"10"`
}

type PageInfo struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}

func Paginate(p *PageQuery) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if p.Page < 1 {
			p.Page = 1
		}

		offset := (p.Page - 1) * p.Size
		return db.Offset(offset).Limit(p.Size)
	}
}
