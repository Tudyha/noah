package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	myio "noah/pkg/io"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"github.com/gorilla/websocket"
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
	client, err := h.clientService.Delete(ctx, GetClientID(ctx))
	if err != nil {
		Fail(ctx, err)
		return
	}

	// 通知客户端退出程序
	if err := session.GetSessionManager().SendCommand(client.SessionID, packet.Command_EXIT); err != nil {
		logger.Info("client logout fail", "err", err)
	}
	Success(ctx, nil)
}

func (h *ClientController) GetClientStat(ctx *gin.Context) {
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

func (h *ClientController) OpenPty(ctx *gin.Context) {
	logger.Info("open pty")
	client, err := h.clientService.GetByID(ctx, GetClientID(ctx))
	if err != nil {
		Fail(ctx, err)
		return
	}
	src, err := session.GetSessionManager().OpenTunnel(client.SessionID, packet.OpenTunnel_PTY, "")
	if err != nil {
		Fail(ctx, err)
		return
	}
	maxMessageSize := 32 * 1024
	upgrader := websocket.Upgrader{
		ReadBufferSize:  maxMessageSize,
		WriteBufferSize: maxMessageSize,
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	target, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		Fail(ctx, err)
		src.Close()
		return
	}
	go func() {
		defer func() {
			src.Close()
			target.Close()
		}()
		tg := &myio.WebSocketReadWriteCloser{Conn: target, MessageType: websocket.TextMessage}
		go io.Copy(tg, src)
		io.Copy(src, tg)
	}()
}
