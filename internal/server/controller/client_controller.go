package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"noah/internal/server/environment"
	"noah/internal/server/gateway"
	"noah/internal/server/middleware/auth"
	"noah/internal/server/middleware/log"
	"noah/internal/server/model"
	"noah/internal/server/service"
	"noah/pkg/enum"
	"noah/pkg/errcode"
	"noah/pkg/mux"
	"noah/pkg/request"
	"os"
	"os/exec"
	"strconv"
	"time"

	myio "noah/pkg/io"

	"noah/pkg/response"

	"noah/pkg/mux/message"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/copier"
	"github.com/samber/do/v2"
)

type ClientController struct {
	clientService  service.IClientService
	gateway        *gateway.Gateway
	authMiddleware *auth.AuthMiddleware
	env            *environment.Environment
}

func NewClientController(i do.Injector) (ClientController, error) {
	return ClientController{
		clientService:  do.MustInvoke[service.IClientService](i),
		gateway:        do.MustInvoke[*gateway.Gateway](i),
		authMiddleware: do.MustInvoke[*auth.AuthMiddleware](i),
		env:            do.MustInvoke[*environment.Environment](i),
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
		Fail(ctx, err)
		return
	}
	var networkList []response.GetClientNetworkInfoRes
	err = json.Unmarshal([]byte(res), &networkList)
	if err != nil {
		log.Error("error", map[string]interface{}{"res": string(res)})
		Fail(ctx, err)
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

func (c ClientController) Connect(ctx *gin.Context) {
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
	hj, ok := ctx.Writer.(http.Hijacker)
	if !ok {
		Fail(ctx, errcode.ErrInternalError)
		return
	}
	conn, _, err := hj.Hijack()
	if err != nil {
		Fail(ctx, errcode.ErrInternalError)
		return
	}

	data, err := json.Marshal(response.Response{
		Code: 0,
		Msg:  "",
		Data: gin.H{"id": id},
	})
	if err != nil {
		Fail(ctx, errcode.ErrInternalError)
		return
	}
	contentLength := len(data)

	response := `HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Content-Length: %d

%s`

	conn.Write([]byte(fmt.Sprintf(response, contentLength, data)))

	// write success response
	go c.gateway.HanderConn(uint32(id), conn)
}

func (c ClientController) GenerateClient(ctx *gin.Context) {
	var req request.GenerateClientReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		Fail(ctx, errcode.ErrInvalidParameter)
		return
	}
	id, err := uuid.NewUUID()
	if err != nil {
		Fail(ctx, err)
		return
	}
	clientBasePath := "client"
	buildStr := `CGO_ENABLED=0 GOOS=%s GOARCH=%s go build -o %s main.go`

	buildCmd := fmt.Sprintf(buildStr, req.Goos, req.Goarch, id.String())

	cmd := exec.Command("sh", "-c", buildCmd)
	cmd.Dir = clientBasePath

	_, err = cmd.CombinedOutput()

	if err != nil {
		Fail(ctx, err)
		return
	}

	defer func() {
		os.Remove(clientBasePath + "/" + id.String())
	}()

	ctx.File(clientBasePath + "/" + id.String())
}

func (c ClientController) GetInstallScript(ctx *gin.Context) {
	r := c.authMiddleware.GenerateTempToken()
	Success(ctx, r.Token)
}
