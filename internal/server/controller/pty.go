package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"noah/internal/server/config"
	"noah/internal/server/service"
	"noah/internal/server/utils"
	"strconv"
)

type PtyController struct{}

func NewPtyController() *PtyController {
	return &PtyController{}
}

// NewPtyChannel 新建pty通道
func (h PtyController) NewPtyChannel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	uintId := uint(id)

	//发送命令，让客户端请求接口建立websocket连接
	channelId := utils.RandString(16)
	_, err := service.GetClientService().SendCommand(uintId, "pty", channelId)
	if err != nil {
		Fail(c, 500, "客户端未上线, shell打开失败")
		return
	}
	//建立与前端的websocket连接
	conn, err := config.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		Fail(c, 500, "Upgrade fail")
		return
	}

	err = service.GetPtyService().NewPtyChannel(channelId, conn)
	if err != nil {
		Fail(c, 500, "NewPtyClient fail")
		return
	}
}

// NewPtyClient 新建pty客户端
func (h PtyController) NewPtyClient(ctx *gin.Context) {
	channelId := ctx.Param("channelId")
	ws, err := config.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	err = service.GetPtyService().NewPtyClient(channelId, ws)
	if err != nil {
		Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
}
