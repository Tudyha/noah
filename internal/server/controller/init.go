package controller

import (
	"errors"
	"net/http"
	"noah/pkg/errcode"
	"noah/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func Init(i do.Injector) {
	do.Provide(i, NewUserController)
	do.Provide(i, NewClientController)
	do.Provide(i, NewAdminController)
	do.Provide(i, NewFileController)
	do.Provide(i, NewTunnelController)
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, response.Response{
		Code: errcode.ErrSuccess.Code,
		Msg:  errcode.ErrSuccess.Msg,
		Data: data,
	})
}

func Fail(c *gin.Context, err error) {
	var appErr *errcode.AppError
	if errors.As(err, &appErr) {
		c.JSON(http.StatusOK, response.Response{
			Code: appErr.Code,
			Msg:  appErr.Msg,
			Data: nil,
		})
	} else {
		c.JSON(http.StatusOK, response.Response{
			Code: errcode.ErrInternalError.Code,
			Msg:  err.Error(),
			Data: nil,
		})
	}
}
