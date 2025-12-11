package model

type WorkSpace struct {
	BaseModel

	Name        string `gorm:"column:name;type:varchar(64);not null;comment:工作空间名称"`
	Description string `gorm:"column:description;type:varchar(255);comment:工作空间描述"`
	Status      int    `gorm:"column:status;type:tinyint;not null;default:1;comment:状态(0-禁用 1-启用)"`
}

func (WorkSpace) TableName() string {
	return "work_space"
}

type WorkSpaceUser struct {
	BaseModel

	SpaceID uint64 `gorm:"column:space_id;type:bigint unsigned;not null;comment:工作空间ID"`
	UserID  uint64 `gorm:"column:user_id;type:bigint unsigned;not null;comment:用户ID"`
	Role    int    `gorm:"column:role;type:tinyint;not null;default:1;comment:角色(1-管理员 2-成员)"`
}

func (WorkSpaceUser) TableName() string {
	return "work_space_user"
}

type WorkSpaceApp struct {
	BaseModel

	Secret      string `gorm:"column:secret;type:varchar(64);not null;comment:密钥"`
	SpaceID     uint64 `gorm:"column:space_id;type:bigint unsigned;not null;comment:工作空间ID"`
	Name        string `gorm:"column:name;type:varchar(64);not null;comment:应用名称"`
	Description string `gorm:"column:description;type:varchar(255);comment:应用描述"`
	Status      int    `gorm:"column:status;type:tinyint;not null;default:1;comment:状态(0-禁用 1-启用)"`
}

func (WorkSpaceApp) TableName() string {
	return "work_space_app"
}
