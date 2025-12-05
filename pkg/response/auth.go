package response

type LoginResponse struct {
	Token string `json:"token" comment:"登录token"`
}
