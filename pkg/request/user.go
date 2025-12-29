package request

import "noah/pkg/enum"

type LoginRequest struct {
	LoginType enum.LoginType `json:"login_type" binding:"required,oneof=1 2"` // 1:验证码登录 2:密码登录
	Username  string         `json:"username" binding:"required"`             // 用户名、手机号或邮箱
	Password  string         `json:"password"`                                // 密码登录时必填
	Code      string         `json:"code"`                                    // 验证码登录时必填
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=4,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Code     string `json:"code" binding:"required"` // 邮箱验证码
}

type SendCodeRequest struct {
	Type     enum.SmsCodeType `json:"type" binding:"required,oneof=1 2 3"` // 1:登录 2:注册 3:邮箱
	Target   string           `json:"target" binding:"required"`           // 手机号或邮箱
}
