package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"noah/internal/server/dao"
	"noah/internal/server/enum"
	"noah/internal/server/gateway"
	"noah/internal/server/middleware"
	"noah/internal/server/middleware/log"
	"noah/internal/server/request"
	"noah/internal/server/response"
	"noah/internal/server/service"
	"noah/pkg/conn"
	"strconv"
	"strings"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/jinzhu/copier"
	"github.com/samber/do/v2"

	"github.com/gin-gonic/gin"
)

type ClientController struct {
	gateway *gateway.Gateway
}

func NewClientController(i do.Injector) *ClientController {
	return &ClientController{
		gateway: do.MustInvoke[*gateway.Gateway](i),
	}
}

// CreateClient 新增客户端
func (c ClientController) CreateClient(ctx *gin.Context) {
	var body request.CreateClientReq
	if err := ctx.ShouldBindJSON(&body); err != nil {
		Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	log.Info("new client", map[string]interface{}{"client": body})

	var client dao.Client
	copier.Copy(&client, body)

	client.RemoteIp = ctx.RemoteIP()
	client.OsType = enum.DetectOS(body.OSName)
	client.LocalIp = body.IPAddress

	id, err := service.GetClientService().Save(client)
	if err != nil {
		Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	Success(ctx, id)
}

// GetClientPage 获取客户端列表
func (c ClientController) GetClientPage(ctx *gin.Context) {
	var req request.ListClientQueryReq
	if err := ctx.BindQuery(&req); err != nil {
		Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	total, clients := service.GetClientService().GetClientPage(req)

	Success(ctx, &request.PageInfo{
		List:  clients,
		Total: total,
	})
}

// GetClient 获取客户端信息
func (c ClientController) GetClient(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	client, err := service.GetClientService().GetClient(uint(id))
	if err != nil {
		Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	Success(ctx, client)
}

// DeleteClient 删除客户端
func (c ClientController) DeleteClient(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	//发送命令让客户端退出
	c.gateway.SendCommand(uint(id), conn.Exit, nil, false)

	//删除客户端
	err := service.GetClientService().Delete(uint(id))
	if err != nil {
		Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	Success(ctx, nil)
}

// NewWsClient 新建客户端ws连接，保持与客户端的长连接
func (c ClientController) NewWsClient(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	err := c.gateway.NewClientWebsocketConn(uint(id), ctx)
	if err != nil {
		Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
}

// SendCommandHandler 执行命令
func (c ClientController) SendCommandHandler(ctx *gin.Context) {
	var form request.SendCommandReq
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if len(strings.TrimSpace(form.Command)) == 0 {
		Fail(ctx, http.StatusBadRequest, "command is empty")
		return
	}

	id := form.ID
	command := &request.CommandRequest{
		Command: form.Command,
	}

	res, err := c.gateway.SendCommand(id, conn.Command, command, true)
	if err != nil {
		Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	Success(ctx, res)
}

// Generate 生成客户端文件
func (c ClientController) Generate(ctx *gin.Context) {
	var req request.ClientGenerateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	filename, err := generate(req)
	if err != nil {
		Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	Success(ctx, filename)
}

// generate 生成客户端文件
// return filename 客户端文件名
func generate(req request.ClientGenerateReq) (string, error) {
	if len(strings.TrimSpace(req.ServerAddr)) == 0 {
		return "", errors.New("serverAddr is empty")
	}

	if len(strings.TrimSpace(req.Port)) == 0 {
		return "", errors.New("port is empty")
	}

	token, err := middleware.GetToken()
	if err != nil {
		return "", err
	}

	filename, err := service.GetClientService().Generate(req.ServerAddr, req.Port, req.OsType, token, req.Filename)
	if err != nil {
		return "", err
	}
	return filename, nil
}

// Update 更新客户端
func (c ClientController) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	//生成最新客户端
	var req request.ClientGenerateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	filename, err := generate(req)
	if err != nil {
		Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	//发送命令让客户端升级
	c.gateway.SendCommand(uint(id), conn.Update, filename, false)
	Success(ctx, "success")
}

// GetClientInfo 获取客户端信息
func (c ClientController) GetClientInfo(ctx *gin.Context) {
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
		startTime = carbon.Parse(start).StdTime()
		endTime = carbon.Parse(end).StdTime()
	}

	clientInfoList, err := service.GetClientService().GetSystemInfo(uint(id), startTime, endTime)
	if err != nil {
		Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	Success(ctx, clientInfoList)
}

func (c ClientController) GetClientProcessList(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	result, err := c.gateway.SendCommand(uint(id), conn.SystemInfo, request.SystemInfoReq{
		SystemInfoType: "process",
		Action:         "list",
		Params:         "",
	}, true)
	if err != nil {
		Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	var processList []response.GetClientProcessRes
	err = json.Unmarshal([]byte(result), &processList)
	if err != nil {
		Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}
	Success(ctx, processList)
}

func (c ClientController) KillClientProcess(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	pid := ctx.Param("pid")
	_, err := c.gateway.SendCommand(uint(id), conn.SystemInfo, request.SystemInfoReq{
		SystemInfoType: "process",
		Action:         "kill",
		Params:         pid,
	}, true)
	if err != nil {
		Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	Success(ctx, "success")
}

func (c ClientController) GetClientNetworkList(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	res, err := c.gateway.SendCommand(uint(id), conn.SystemInfo, request.SystemInfoReq{
		SystemInfoType: "net",
		Action:         "list",
		Params:         "",
	}, true)
	if err != nil {
		Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	var networkList []response.GetClientNetworkInfoRes
	err = json.Unmarshal([]byte(res), &networkList)
	if err != nil {
		Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}
	Success(ctx, networkList)
}

func (c ClientController) GetClientDockerContainerList(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	res, err := c.gateway.SendCommand(uint(id), conn.SystemInfo, request.SystemInfoReq{
		SystemInfoType: "docker",
		Action:         "containerList",
		Params:         "",
	}, true)
	if err != nil {
		Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	var containerList []response.GetClientDockerContainerRes
	err = json.Unmarshal([]byte(res), &containerList)
	if err != nil {
		Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}
	Success(ctx, containerList)
}

// Health 心跳检测
func (c ClientController) Health(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var systemInfo request.CreateSystemInfoReq
	if err := ctx.ShouldBindJSON(&systemInfo); err != nil {
		Fail(ctx, 400, "")
		return
	}
	service.GetClientService().UpdateStatus(uint(id), enum.DEVICE_ONLINE)

	if systemInfo.CpuUsage != 0 {
		err := service.GetClientService().SaveSystemInfo(uint(id), systemInfo)
		if err != nil {
			Fail(ctx, 400, "")
			return
		}
	}
	Success(ctx, nil)
}
