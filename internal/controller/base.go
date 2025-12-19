package controller

import (
	"errors"
	"net/http"
	"noah/pkg/constant"
	"noah/pkg/errcode"
	"noah/pkg/response"
	"noah/pkg/utils"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	once sync.Once

	authController   *AuthController
	userController   *UserController
	clientController *ClientController
)

func Init() error {
	once.Do(func() {
		authController = newAuthController()
		userController = newUserController()
		clientController = newClientController()
	})
	return nil
}

func GetAuthController() *AuthController {
	return authController
}

func GetUserController() *UserController {
	return userController
}

func GetClientController() *ClientController {
	return clientController
}

func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, response.Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

func Fail(ctx *gin.Context, err error) {
	FailWithMsg(ctx, err, "")
}

func FailWithMsg(ctx *gin.Context, err error, msg string) {
	code, msg := getErrorMsg(err)

	ctx.JSON(http.StatusOK, response.Response{
		Code: code,
		Msg:  msg,
	})
}

func getErrorMsg(err error) (code int, msg string) {
	var appErr *errcode.AppError
	if errors.As(err, &appErr) {
		if msg == "" {
			msg = appErr.Msg
		}
		code = appErr.Code
	} else {
		if msg == "" {
			msg = errcode.ErrInternalServer.Msg
		}
		code = errcode.ErrInternalServer.Code
	}
	return
}

func GetUserId(ctx *gin.Context) uint64 {
	return ctx.GetUint64(constant.HttpHeaderUserIDKey)
}

func GetAppID(ctx *gin.Context) uint64 {
	id, err := utils.StringToUint64(ctx.Request.Header.Get(constant.HttpHeaderAppIDKey))
	if err != nil {
		return 0
	}
	return id
}

func GetClientID(ctx *gin.Context) uint64 {
	v, ok := ctx.Params.Get("client_id")
	if !ok {
		return 0
	}
	clientID, err := utils.StringToUint64(v)
	if err != nil {
		return 0
	}
	return clientID
}
