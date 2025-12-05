package model

import (
	"time"
)

type BaseModel struct {
	ID        uint64     `gorm:"primary_key;auto_increment;comment:主键ID" json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at;type:datetime;not null;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:datetime;not null;comment:更新时间" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime;index;comment:删除时间" json:"deleted_at"`
}
