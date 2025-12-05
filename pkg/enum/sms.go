package enum

type SmsCodeType uint

const (
	SmsCodeTypeLogin SmsCodeType = iota + 1 // 登录验证码
)
