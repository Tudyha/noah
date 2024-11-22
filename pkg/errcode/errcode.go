package errcode

const (
	ErrCodeSuccess          = 0
	ErrCodeInternalError    = 500000
	ErrCodeUserNotFound     = 501000
	ErrCodePasswordNotMatch = 501001
	ErrCodePermissionDenied = 501002
)

const (
	ErrCodeInvalidParameter = 400
	ErrCodeTokenExpired     = 401
)

var (
	ErrSuccess          = &AppError{Code: ErrCodeSuccess, Msg: "成功"}
	ErrInternalError    = &AppError{Code: ErrCodeInternalError, Msg: "内部错误"}
	ErrInvalidParameter = &AppError{Code: ErrCodeInvalidParameter, Msg: "无效参数"}
	ErrUserNotFound     = &AppError{Code: ErrCodeUserNotFound, Msg: "用户不存在"}
	ErrPasswordNotMatch = &AppError{Code: ErrCodePasswordNotMatch, Msg: "用户名/密码不正确"}
	ErrTokenExpired     = &AppError{Code: ErrCodeTokenExpired, Msg: "token过期"}
	ErrPermissionDenied = &AppError{Code: ErrCodePermissionDenied, Msg: "权限不足"}
)

type AppError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *AppError) Error() string {
	return e.Msg
}
