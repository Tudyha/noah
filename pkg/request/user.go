package request

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdatePasswordReq struct {
	Password string `json:"password" binding:"required"`
}
