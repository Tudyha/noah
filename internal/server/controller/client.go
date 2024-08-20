package controller

import (
	"net/http"
	"noah/internal/server/config"
	"noah/internal/server/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ClientController struct{}

func NewClientController() *ClientController {
	return &ClientController{}
}

func (h *ClientController) NewClient(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ws, err := config.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = service.GetClientService().AddConnection(uint(id), ws)
	if err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *ClientController) NewPtyClient(c *gin.Context) {
	channelId := c.Param("channelId")
	ws, err := config.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = service.GetPtyService().AddPtyConnection(channelId, ws)
	if err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
}

type SendCommandRequestForm struct {
	ID        uint   `json:"id" binding:"required"`
	Command   string `json:"command" binding:"required"`
	Parameter string `json:"parameter"`
}

func (h *ClientController) SendCommandHandler(c *gin.Context) {
	var form SendCommandRequestForm
	if err := c.ShouldBindJSON(&form); err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if len(strings.TrimSpace(form.Command)) == 0 {
		Fail(c, http.StatusBadRequest, "command is empty")
		return
	}

	id := form.ID

	res, err := service.GetClientService().SendCommand(id, form.Command, form.Parameter)
	if err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, res)
}

type ClientGenerateReq struct {
	ServerAddr string `json:"serverAddr" binding:"required"`
	Port       string `json:"port" binding:"required"`
	OsType     int8   `json:"osType" binding:"required"`
	Filename   string `json:"filename"`
}

func (h *ClientController) Generate(c *gin.Context) {
	var req ClientGenerateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	if len(strings.TrimSpace(req.ServerAddr)) == 0 {
		Fail(c, http.StatusBadRequest, "serverAddr is empty")
		return
	}

	if len(strings.TrimSpace(req.Port)) == 0 {
		Fail(c, http.StatusBadRequest, "port is empty")
		return
	}

	filename, err := service.GetClientService().Generate(req.ServerAddr, req.Port, req.OsType, req.Filename)
	if err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	Success(c, filename)
}

func (h *ClientController) Download(c *gin.Context) {
	filename := c.Param("filename")

	c.File("temp/" + filename)
}
