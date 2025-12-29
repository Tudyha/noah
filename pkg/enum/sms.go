package enum

type SmsCodeType uint

const (
	SmsCodeTypeLogin    SmsCodeType = iota + 1 // 登录验证码
	SmsCodeTypeRegister                        // 注册验证码
	SmsCodeTypeEmail                           // 邮箱验证码
)
