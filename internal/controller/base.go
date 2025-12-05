package controller

import (
	"errors"
	"net/http"
	"noah/pkg/constant"
	"noah/pkg/errcode"
	"noah/pkg/response"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	once sync.Once

	authController *AuthController
	userController *UserController
)

func Init() error {
	once.Do(func() {
		authController = newAuthController()
		userController = newUserController()
	})
	return nil
}

func GetAuthController() *AuthController {
	return authController
}

func GetUserController() *UserController {
	return userController
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
