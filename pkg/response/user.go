package response

type UserResponse struct {
	ID       uint64 `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`

	WorkSpaceList []*WorkSpaceResponse `json:"work_space_list"` // 工作空间列表
}
