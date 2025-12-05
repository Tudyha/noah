package model

type WorkSpace struct {
	BaseModel

	Name        string `gorm:"column:name;type:varchar(64);not null;comment:工作空间名称" json:"name"`
	Description string `gorm:"column:description;type:varchar(255);comment:工作空间描述" json:"description"`
	Status      int    `gorm:"column:status;type:tinyint;not null;default:1;comment:状态(0-禁用 1-启用)" json:"status"`
}

func (WorkSpace) TableName() string {
	return "work_space"
}

type WorkSpaceUser struct {
	BaseModel

	SpaceID uint64 `gorm:"column:space_id;type:bigint unsigned;not null;comment:工作空间ID" json:"space_id"`
	UserID  uint64 `gorm:"column:user_id;type:bigint unsigned;not null;comment:用户ID" json:"user_id"`
	Role    int    `gorm:"column:role;type:tinyint;not null;default:1;comment:角色(1-管理员 2-成员)" json:"role"`
}

func (WorkSpaceUser) TableName() string {
	return "work_space_user"
}

type WorkSpaceApp struct {
	BaseModel

	SpaceID     uint64 `gorm:"column:space_id;type:bigint unsigned;not null;comment:工作空间ID" json:"space_id"`
	Name        string `gorm:"column:name;type:varchar(64);not null;comment:应用名称" json:"name"`
	Description string `gorm:"column:description;type:varchar(255);comment:应用描述" json:"description"`
	Status      int    `gorm:"column:status;type:tinyint;not null;default:1;comment:状态(0-禁用 1-启用)" json:"status"`
}

func (WorkSpaceApp) TableName() string {
	return "work_space_app"
}
