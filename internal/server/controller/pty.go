package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"noah/internal/server/config"
	"noah/internal/server/service"
	"noah/internal/server/utils"
	"strconv"
)

type PtyController struct{}

func NewPtyController() *PtyController {
	return &PtyController{}
}

func (h *PtyController) WebSocket(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	uintId := uint(id)

	//先建立客户端websocket连接
	channelId := utils.RandString(16)
	_, err := service.GetClientService().SendCommand(uintId, "pty", channelId)
	if err != nil {
		log.Println("客户端websocket连接失败", err)
		Fail(c, 500, "客户端未上线, shell打开失败")
		return
	}
	//建立与前端的websocket连接
	conn, err := config.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrader.Upgrade:", err)
		Fail(c, 500, "Upgrade fail")
		return
	}

	err = service.GetPtyService().NewPtyClient(channelId, conn)
	if err != nil {
		log.Println("NewPtyClient:", err)
		Fail(c, 500, "NewPtyClient fail")
		return
	}
}
