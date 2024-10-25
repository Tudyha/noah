package controller

import (
	"net/http"
	"noah/internal/server/enum"
	"noah/internal/server/gateway"
	"noah/internal/server/middleware/log"
	"noah/internal/server/request"
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
	conn, err := gateway.UpgradeWebsocket(c)
	if err != nil {
		Fail(c, 500, "Upgrade fail")
		return
	}

	err = service.GetChannelService().NewChannel(uintId, request.CreateChannelReq{
		ChannelType: enum.Pty,
	}, conn)
	if err != nil {
		log.Error("open pty fail", map[string]interface{}{"err": err.Error()})

		// 断开与前端的websocket连接
		conn.Close()
		return
	}
}

func (h ChannelController) NewChannel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	uintId := uint(id)

	var req request.CreateChannelReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	err = service.GetChannelService().NewChannel(uintId, req, nil)
	if err != nil {
		Fail(c, 500, err.Error())
		return
	}

	Success(c, "success")
}

func (h ChannelController) GetChannelList(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	uintId := uint(id)

	res, err := service.GetChannelService().GetChannelList(uintId)
	if err != nil {
		Fail(c, 500, err.Error())
		return
	}
	Success(c, res)

}

func (h ChannelController) DeleteChannel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("channelId"))
	uintId := uint(id)

	err := service.GetChannelService().DeleteChannel(uintId)
	if err != nil {
		Fail(c, 500, err.Error())
		return
	}
	Success(c, "success")
}
