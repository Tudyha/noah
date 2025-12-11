package response

import "time"

type WorkSpaceResponse struct {
	ID          uint64     `json:"id"`          // 工作空间ID
	Name        string     `json:"name"`        // 工作空间名称
	Description string     `json:"description"` // 工作空间描述
	Status      int        `json:"status"`      // 工作空间状态
	CreatedAt   *time.Time `json:"created_at"`  // 创建时间

	AppList []*WorkSpaceAppResponse `json:"app_list"` // 应用列表
}

type WorkSpaceAppResponse struct {
	ID          uint64     `json:"id"`          // 应用ID
	Name        string     `json:"name"`        // 应用名称
	Description string     `json:"description"` // 应用描述
	Status      int        `json:"status"`      // 应用状态
	CreatedAt   *time.Time `json:"created_at"`  // 创建时间
}
