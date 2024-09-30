package controller

import (
	"noah/internal/server/gateway"
)

type Controller struct {
	ClientController  *ClientController
	clientController  *ClientController
	channelController *ChannelController
	userController    *UserController
	fileController    *FileController
	shellController   *ShellController
	adminController   *AdminController
}

func NewController(gateway *gateway.Gateway) *Controller {
	return &Controller{
		clientController:  NewClientController(gateway),
		shellController:   NewShellController(),
		channelController: NewChannelController(),
		userController:    NewUserController(),
		fileController:    NewFileController(gateway),
		adminController:   NewAdminController(),
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

func (c Controller) GetAdminController() *AdminController {
	return c.adminController
}
