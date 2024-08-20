package dto

type LoginReq struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type UserInfo struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	Token    string `json:"token"`
	Expire   int64  `json:"expire"`
	Avatar   string `json:"avatar"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   int8   `json:"status"`
	Created  int64  `json:"created"`
	Updated  int64  `json:"updated"`
}
