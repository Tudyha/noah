package request

import "noah/pkg/enum"

type LoginRequest struct {
	LoginType enum.LoginType `json:"loginType" binding:"required,oneof=1 2"` // 1:验证码登录 2:密码登录
	Username  string         `json:"username" binding:"required"`            // 用户名或者手机号
	Password  string         `json:"password"`                               // 密码登录时必填
	Code      string         `json:"code"`                                   // 验证码登录时必填
}
