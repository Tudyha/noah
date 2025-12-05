package errcode

const (
	// 认证相关错误 2000-2999
	AUTH_ERROR = 2000 + iota
)

const (
	// 用户相关错误 3000-3999
	USER_ERROR = 3000 + iota
)

var userErrorMessages = map[int]string{

	AUTH_ERROR: "认证相关错误",

	USER_ERROR: "用户相关错误",
}

var (
	// 认证相关错误
	ErrAuth        = &AppError{Code: AUTH_ERROR, Msg: userErrorMessages[AUTH_ERROR]}
	ErrLoginFailed = &AppError{Code: AUTH_ERROR + 1, Msg: "登录失败，用户名或密码错误"}

	// 用户相关错误
	ErrUser               = &AppError{Code: USER_ERROR, Msg: userErrorMessages[USER_ERROR]}
	ErrUserNotFound       = &AppError{Code: USER_ERROR + 1, Msg: "用户未找到"}
	ErrUserDisabled       = &AppError{Code: USER_ERROR + 2, Msg: "用户已被禁用"}
	ErrUserNotSetPassword = &AppError{Code: USER_ERROR + 3, Msg: "用户未设置密码"}
)
