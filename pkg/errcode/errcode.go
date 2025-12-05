package errcode

const (
	INVALID_PARAMS        = 400
	UNAUTHORIZED          = 401
	FORBIDDEN             = 403
	NOT_FOUND             = 404
	METHOD_NOT_ALLOWED    = 405
	INTERNAL_SERVER_ERROR = 500
	SERVICE_UNAVAILABLE   = 503
)

var commondErrorMessages = map[int]string{

	INVALID_PARAMS:        "参数错误",
	UNAUTHORIZED:          "未授权",
	FORBIDDEN:             "无权限",
	NOT_FOUND:             "资源不存在",
	METHOD_NOT_ALLOWED:    "方法不允许",
	INTERNAL_SERVER_ERROR: "服务器内部错误",
	SERVICE_UNAVAILABLE:   "服务不可用",
}

var (
	// 通用错误
	ErrInvalidParams      = &AppError{Code: INVALID_PARAMS, Msg: commondErrorMessages[INVALID_PARAMS]}
	ErrUnauthorized       = &AppError{Code: UNAUTHORIZED, Msg: commondErrorMessages[UNAUTHORIZED]}
	ErrForbidden          = &AppError{Code: FORBIDDEN, Msg: commondErrorMessages[FORBIDDEN]}
	ErrNotFound           = &AppError{Code: NOT_FOUND, Msg: commondErrorMessages[NOT_FOUND]}
	ErrMethodNotAllowed   = &AppError{Code: METHOD_NOT_ALLOWED, Msg: commondErrorMessages[METHOD_NOT_ALLOWED]}
	ErrInternalServer     = &AppError{Code: INTERNAL_SERVER_ERROR, Msg: commondErrorMessages[INTERNAL_SERVER_ERROR]}
	ErrServiceUnavailable = &AppError{Code: SERVICE_UNAVAILABLE, Msg: commondErrorMessages[SERVICE_UNAVAILABLE]}
)

// AppError 应用错误
type AppError struct {
	Code int
	Msg  string
}

func (e *AppError) Error() string {
	return e.Msg
}
