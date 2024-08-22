package controller

import (
	"noah/internal/server/enum"
	"noah/internal/server/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	deviceController *DeviceController
	clientController *ClientController
	shellController  *ShellController
	ptyController    *PtyController
	userController   *UserController
	fileController   *FileController
}

func NewController() *Controller {
	return &Controller{
		deviceController: NewDeviceController(),
		clientController: NewClientController(),
		shellController:  NewShellController(),
		ptyController:    NewPtyController(),
		userController:   NewUserController(),
		fileController:   NewFileController(),
	}
}

func (c *Controller) GetDeviceController() *DeviceController {
	return c.deviceController
}

func (c *Controller) GetClientController() *ClientController {
	return c.clientController
}

func (c *Controller) GetShellController() *ShellController {
	return c.shellController
}

func (c *Controller) GetPtyController() *PtyController {
	return c.ptyController
}

func (c *Controller) GetUserController() *UserController {
	return c.userController
}

func (c *Controller) GetFileController() *FileController {
	return c.fileController
}

func (h *Controller) Health(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	service.GetDeviceService().UpdateStatus(uint(id), enum.DEVICE_ONLINE)
	Success(c, nil)
}
