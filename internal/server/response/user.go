package response

type GetUserInfoRes struct {
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
