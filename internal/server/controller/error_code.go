package controller

type ErrorCode struct {
	Code int
	Msg  string
}

var (
	ServerUnknowedError = ErrorCode{50000, "server unknowed error"}
)

func (e ErrorCode) Error() string {
	return e.Msg
}

func (e ErrorCode) GetCode() int {
	return e.Code
}
