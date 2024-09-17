package controller

import (
	"fmt"
	"noah/internal/server/config"
	"noah/internal/server/enum"
	"noah/internal/server/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChannelController struct{}

func NewChannelController() *ChannelController {
	return &ChannelController{}
}

// NewPtyChannel 新建pty通道
func (h ChannelController) NewPtyChannel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	uintId := uint(id)

	//建立与前端的websocket连接
	conn, err := config.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		Fail(c, 500, "Upgrade fail")
		return
	}

	err = service.GetChannelService().NewChannel(uintId, enum.Pty, conn, "", "")
	if err != nil {
		Fail(c, 500, "NewChannel fail")
		return
	}

	if err != nil {
		Fail(c, 500, "NewPtyClient fail")
		return
	}
}

type CreateChannelReq struct {
	ChannelType enum.ChannelType `json:"channelType" binding:"required"`
	ServerPort  string           `json:"serverPort"`
	ClientIp    string           `json:"clientIp"`
	ClientPort  string           `json:"clientPort"`
}

func (h ChannelController) NewChannel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	uintId := uint(id)

	var req CreateChannelReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Fail(c, 500, "参数错误")
		return
	}

	err = service.GetChannelService().NewChannel(uintId, req.ChannelType, nil, req.ServerPort, fmt.Sprintf("%s:%s", req.ClientIp, req.ClientPort))
	if err != nil {
		Fail(c, 500, err.Error())
		return
	}

	Success(c, "success")
}
