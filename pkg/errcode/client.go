package errcode

const (
	// 客户端相关错误 4000-4999
	CLIENT_ERROR = 4000 + iota
	CLIENT_NOT_FOUND
	CLIENT_DISCONNECT
)

var (
	ErrClientDisconnect = &AppError{Code: CLIENT_DISCONNECT, Msg: "客户端未连接"}
)
