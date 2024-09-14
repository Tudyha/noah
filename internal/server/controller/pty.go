package controller

import (
	"noah/internal/server/config"
	"noah/internal/server/enum"
	"noah/internal/server/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PtyController struct{}

func NewPtyController() *PtyController {
	return &PtyController{}
}

// NewPtyChannel 新建pty通道
func (h PtyController) NewPtyChannel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	uintId := uint(id)

	//建立与前端的websocket连接
	conn, err := config.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		Fail(c, 500, "Upgrade fail")
		return
	}

	_, err = service.GetChannelService().NewChannel(uintId, enum.Pty, conn, "")
	if err != nil {
		Fail(c, 500, "NewChannel fail")
		return
	}

	if err != nil {
		Fail(c, 500, "NewPtyClient fail")
		return
	}
}

// NewPtyClient 新建pty客户端
func (h PtyController) NewPtyClient(ctx *gin.Context) {
	//channelId := ctx.Param("channelId")
	//ws, err := config.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	//if err != nil {
	//	Fail(ctx, http.StatusInternalServerError, err.Error())
	//	return
	//}
	//
	//err = service.GetChannelService().ClientConnect(channelId, ws)
	//if err != nil {
	//	Fail(ctx, http.StatusInternalServerError, err.Error())
	//	return
	//}
}

type CreateChannelReq struct {
	ChannelType enum.ChannelType `json:"channelType" binding:"required"`
	ServerPort  string           `json:"serverPort"`
	ClientIp    string           `json:"clientIp"`
	ClientPort  string           `json:"clientPort"`
}

func (h PtyController) NewChannel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	uintId := uint(id)

	var req CreateChannelReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Fail(c, 500, "参数错误")
		return
	}

	_, err = service.GetChannelService().NewChannel(uintId, req.ChannelType, nil, req.ServerPort)
	if err != nil {
		Fail(c, 500, err.Error())
		return
	}

	Success(c, "success")
}

func (h PtyController) ChannelClientConnect(ctx *gin.Context) {
	//channelId := ctx.Param("channelId")
	//ws, err := config.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	//if err != nil {
	//	Fail(ctx, http.StatusInternalServerError, err.Error())
	//	return
	//}
	//
	//err = service.GetChannelService().ClientConnect(channelId, ws)
	//if err != nil {
	//	Fail(ctx, http.StatusInternalServerError, err.Error())
	//	return
	//}
}
