package controller

import (
	"noah/internal/server/enum"
	"noah/internal/server/gateway"
	"noah/internal/server/request"
	"noah/internal/server/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	ClientController  *ClientController
	clientController  *ClientController
	channelController *ChannelController
	userController    *UserController
	fileController    *FileController
	shellController   *ShellController
}

func NewController(gateway *gateway.Gateway) *Controller {
	return &Controller{
		clientController:  NewClientController(gateway),
		shellController:   NewShellController(),
		channelController: NewChannelController(),
		userController:    NewUserController(),
		fileController:    NewFileController(gateway),
	}
}

func (c Controller) GetClientController() *ClientController {
	return c.clientController
}

func (c Controller) GetShellController() *ShellController {
	return c.shellController
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

// Health 心跳检测
func (c Controller) Health(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var systemInfo request.CreateSystemInfoReq
	if err := ctx.ShouldBindJSON(&systemInfo); err != nil {
		Fail(ctx, 400, "")
		return
	}
	service.GetClientService().UpdateStatus(uint(id), enum.DEVICE_ONLINE)

	if systemInfo.CpuUsage != 0 {
		err := service.GetClientService().SaveSystemInfo(uint(id), systemInfo)
		if err != nil {
			Fail(ctx, 400, "")
			return
		}
	}
	Success(ctx, nil)
}
