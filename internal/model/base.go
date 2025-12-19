package model

import (
	"time"
)

type BaseModel struct {
	ID        uint64    `gorm:"primary_key;auto_increment;comment:主键ID"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;comment:创建时间"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;comment:更新时间"`
}
