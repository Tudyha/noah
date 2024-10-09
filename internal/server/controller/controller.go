package controller

import "github.com/samber/do/v2"

type Controller struct {
	ClientController  *ClientController
	clientController  *ClientController
	channelController *ChannelController
	userController    *UserController
	fileController    *FileController
	shellController   *ShellController
	adminController   *AdminController
}

func NewController(i do.Injector) *Controller {
	return &Controller{
		clientController:  NewClientController(i),
		shellController:   NewShellController(),
		channelController: NewChannelController(),
		userController:    NewUserController(),
		fileController:    NewFileController(i),
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
