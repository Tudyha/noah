package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"noah/internal/server/gateway"
	"noah/internal/server/model"
	"noah/internal/server/service"
	"noah/pkg/enum"
	"noah/pkg/errcode"
	"noah/pkg/mux"
	"noah/pkg/request"
	"strconv"
	"time"

	myio "noah/pkg/io"

	"noah/pkg/response"

	"noah/pkg/mux/message"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/copier"
	"github.com/samber/do/v2"
)

type ClientController struct {
	clientService service.IClientService
	gateway       *gateway.Gateway
}

func NewClientController(i do.Injector) (ClientController, error) {
	return ClientController{
		clientService: do.MustInvoke[service.IClientService](i),
		gateway:       do.MustInvoke[*gateway.Gateway](i),
	}, nil
}

// CreateClient 新增客户端
func (c ClientController) CreateClient(ctx *gin.Context) {
	var body request.CreateClientReq
	if err := ctx.ShouldBindJSON(&body); err != nil {
		Fail(ctx, errcode.ErrInvalidParameter)
		return
	}

	var client model.Client
	copier.Copy(&client, body)

	client.RemoteIp = ctx.RemoteIP()
	client.OsType = enum.DetectOS(body.OSName)
	client.LocalIp = body.IPAddress

	id, err := c.clientService.Save(client)
	if err != nil {
		Fail(ctx, err)
		return
	}
	Success(ctx, id)
}

// GetClientPage 获取客户端列表
func (c ClientController) GetClientPage(ctx *gin.Context) {
	var req request.ListClientQueryReq
	if err := ctx.BindQuery(&req); err != nil {
		Fail(ctx, errcode.ErrInvalidParameter)
		return
	}

	total, clients := c.clientService.GetClientPage(req)

	Success(ctx, &request.PageInfo{
		List:  clients,
		Total: total,
	})
}

// GetClient 获取客户端信息
func (c ClientController) GetClient(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	client, err := c.clientService.GetClient(uint(id))
	if err != nil {
		Fail(ctx, errcode.ErrInternalError)
		return
	}

	Success(ctx, client)
}

// DeleteClient 删除客户端
func (c ClientController) DeleteClient(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	//发送命令让客户端退出
	c.gateway.SendCommand(uint(id), enum.Exit, nil, false)

	//删除客户端
	err := c.clientService.Delete(uint(id))
	if err != nil {
		Fail(ctx, errcode.ErrInternalError)
		return
	}

	Success(ctx, nil)
}

// NewPtyChannel 新建pty通道
func (c ClientController) OpenPty(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	//建立与前端的websocket连接
	maxMessageSize := 32 * 1024
	upgrader := websocket.Upgrader{
		ReadBufferSize:  maxMessageSize,
		WriteBufferSize: maxMessageSize,
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	srcConn, err := c.gateway.NewClientConn(uint32(id), "pty", "")
	if err != nil {
		Fail(ctx, errcode.ErrInternalError)
		return
	}

	targetConn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		Fail(ctx, errcode.ErrInternalError)
		return
	}

	go copy(srcConn, &myio.WebSocketReaderWriterCloser{Conn: targetConn})
}

func copy(src, target io.ReadWriteCloser) {
	defer src.Close()
	defer target.Close()
	src.(*mux.Conn).Copy(target)
}

// GetClientInfo 获取客户端信息
func (c ClientController) GetClientStat(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	start := ctx.Query("start")
	end := ctx.Query("end")
	var startTime, endTime time.Time
	if start == "" || end == "" {
		//获取当前时间
		endTime = time.Now()
		//获取5分钟前时间
		startTime = endTime.Add(-5 * time.Minute)
	} else {
		startTime = carbon.Parse(start).ToStdTime()
		endTime = carbon.Parse(end).ToStdTime()
	}

	clientInfoList, err := c.clientService.GetClientStat(uint(id), startTime, endTime)
	if err != nil {
		Fail(ctx, errcode.ErrInternalError)
		return
	}
	Success(ctx, clientInfoList)
}

func (c ClientController) GetClientProcessList(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	result, err := c.gateway.SendCommand(uint(id), enum.SystemInfo, message.SystemInfoReq{
		SystemInfoType: "process",
		Action:         "list",
		Params:         "",
	}, true)
	if err != nil {
		Fail(ctx, errcode.ErrInternalError)
		return
	}

	var processList []response.GetClientProcessRes
	err = json.Unmarshal([]byte(result), &processList)
	if err != nil {
		Fail(ctx, errcode.ErrInternalError)
		return
	}
	Success(ctx, processList)
}

func (c ClientController) KillClientProcess(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	pid := ctx.Param("pid")
	_, err := c.gateway.SendCommand(uint(id), enum.SystemInfo, message.SystemInfoReq{
		SystemInfoType: "process",
		Action:         "kill",
		Params:         pid,
	}, true)
	if err != nil {
		Fail(ctx, errcode.ErrInternalError)
		return
	}
	Success(ctx, "success")
}

func (c ClientController) GetClientNetworkList(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	res, err := c.gateway.SendCommand(uint(id), enum.SystemInfo, message.SystemInfoReq{
		SystemInfoType: "net",
		Action:         "list",
		Params:         "",
	}, true)
	if err != nil {
		Fail(ctx, errcode.ErrInternalError)
		return
	}
	var networkList []response.GetClientNetworkInfoRes
	err = json.Unmarshal([]byte(res), &networkList)
	if err != nil {
		Fail(ctx, errcode.ErrInternalError)
		return
	}
	Success(ctx, networkList)
}

func (c ClientController) GetClientDockerContainerList(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	res, err := c.gateway.SendCommand(uint(id), enum.SystemInfo, message.SystemInfoReq{
		SystemInfoType: "docker",
		Action:         "containerList",
		Params:         "",
	}, true)
	if err != nil {
		Fail(ctx, errcode.ErrInternalError)
		return
	}
	var containerList []response.GetClientDockerContainerRes
	err = json.Unmarshal([]byte(res), &containerList)
	if err != nil {
		Fail(ctx, errcode.ErrInternalError)
		return
	}
	Success(ctx, containerList)
}
