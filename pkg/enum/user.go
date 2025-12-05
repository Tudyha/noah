package enum

type LoginType int

const (
	LoginTypeCode     = iota + 1 // 验证码登录
	LoginTypePassword            // 密码登录
)
