package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type Controller struct {
	ClientController  *ClientController
	clientController  *ClientController
	channelController *ChannelController
	userController    *UserController
	fileController    *FileController
	adminController   *AdminController
}

func NewController(i do.Injector) *Controller {
	return &Controller{
		clientController:  NewClientController(i),
		channelController: NewChannelController(),
		userController:    NewUserController(),
		fileController:    NewFileController(i),
		adminController:   NewAdminController(),
	}
}

func (c Controller) GetClientController() *ClientController {
	return c.clientController
}

func (c Controller) GetChannelController() *ChannelController {
	return c.channelController
}

func (c Controller) GetUserController() *UserController {
	return c.userController
}

func (c Controller) GetFileController() *FileController {
	return c.fileController
}

func (c Controller) GetAdminController() *AdminController {
	return c.adminController
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
