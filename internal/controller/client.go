package controller

import (
	"encoding/json"
	"fmt"
	"noah/internal/service"
	"noah/internal/session"
	"noah/pkg/config"
	"noah/pkg/errcode"
	"noah/pkg/logger"
	"noah/pkg/packet"
	"noah/pkg/request"
	"noah/pkg/response"
	"noah/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"github.com/jinzhu/copier"
)

type ClientController struct {
	clientService service.ClientService
	workService   service.WorkService
}

func newClientController() *ClientController {
	return &ClientController{
		clientService: service.GetClientService(),
		workService:   service.GetWorkService(),
	}
}

func (h *ClientController) GetClientPage(ctx *gin.Context) {
	appID := GetAppID(ctx)
	if appID == 0 {
		Fail(ctx, errcode.ErrInvalidParams)
		return
	}

	var req request.ClientQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		FailWithMsg(ctx, errcode.ErrInvalidParams, err.Error())
		return
	}
	res, err := h.clientService.GetPage(ctx, appID, req)
	if err != nil {
		Fail(ctx, err)
		return
	}
	Success(ctx, res)
}

func (h *ClientController) GetClientBind(ctx *gin.Context) {
	appID := GetAppID(ctx)
	if appID == 0 {
		Fail(ctx, errcode.ErrInvalidParams)
		return
	}

	app, err := h.workService.GetAppByAppID(ctx, appID)
	if err != nil {
		Fail(ctx, err)
		return
	}

	cfg := config.Get()

	clientConfig := config.ClientConfig{
		ServerAddr: cfg.Server.TCP.Addr,
		AppId:      appID,
		AppSecret:  app.Secret,
	}

	data, err := json.Marshal(clientConfig)
	if err != nil {
		Fail(ctx, err)
		return
	}

	c := utils.Base64Encode(data)

	sh := "curl -s http://%s/file/%s -o /tmp/noah-cli && chmod +x /tmp/noah-cli && /tmp/noah-cli run -c %s"
	Success(ctx, response.ClientBindResponse{
		MacBind: fmt.Sprintf(sh, cfg.Server.HTTP.Addr, "noah-mac", c),
	})
}

func (h *ClientController) DeleteClient(ctx *gin.Context) {
	appID := GetAppID(ctx)
	if appID == 0 {
		Fail(ctx, errcode.ErrInvalidParams)
		return
	}

	client, err := h.clientService.Delete(ctx, GetClientID(ctx))
	if err != nil {
		Fail(ctx, err)
		return
	}

	// 通知客户端退出程序
	if err := session.GetSessionManager().SendProtoMessage(client.ConnID, packet.MessageType_Logout, &packet.Logout{}); err != nil {
		logger.Info("client logout fail", "err", err)
	}
	Success(ctx, nil)
}

func (h *ClientController) GetClientStat(ctx *gin.Context) {
	appID := GetAppID(ctx)
	if appID == 0 {
		Fail(ctx, errcode.ErrInvalidParams)
		return
	}
	start := ctx.Query("start")
	end := ctx.Query("end")
	var startTime, endTime time.Time
	if start == "" || end == "" {
		//获取当前时间
		endTime = time.Now()
		//获取5分钟前时间
		startTime = endTime.Add(-50 * time.Minute)
	} else {
		startTime = carbon.Parse(start).ToStdTime()
		endTime = carbon.Parse(end).ToStdTime()
	}

	list, err := h.clientService.GetClientStat(ctx, GetClientID(ctx), startTime, endTime)
	if err != nil {
		Fail(ctx, err)
		return
	}
	var clientInfoList []response.ClientStatResponse
	if err := copier.CopyWithOption(&clientInfoList, list, copier.Option{
		Converters: []copier.TypeConverter{{
			SrcType: "",
			DstType: []response.DiskUsageStat{},
			Fn: func(src interface{}) (interface{}, error) {
				var dst []response.DiskUsageStat
				err := json.Unmarshal([]byte(src.(string)), &dst)
				return dst, err
			},
		}},
	}); err != nil {
		Fail(ctx, err)
		return
	}
	Success(ctx, clientInfoList)
}
