package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
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
	"os"
	"strings"
	"time"

	myio "noah/pkg/io"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/copier"
)

type AgentController struct {
	agentService service.AgentService
	workService  service.WorkService
}

func newAgentController() *AgentController {
	return &AgentController{
		agentService: service.GetAgentService(),
		workService:  service.GetWorkService(),
	}
}

func (h *AgentController) GetAgentPage(ctx *gin.Context) {
	appID := GetAppID(ctx)
	if appID == 0 {
		Fail(ctx, errcode.ErrInvalidParams)
		return
	}

	var req request.AgentQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		FailWithMsg(ctx, errcode.ErrInvalidParams, err.Error())
		return
	}
	res, err := h.agentService.GetPage(ctx, appID, req)
	if err != nil {
		Fail(ctx, err)
		return
	}
	Success(ctx, res)
}

func (h *AgentController) GetAgent(ctx *gin.Context) {
	var req request.AgentQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		FailWithMsg(ctx, errcode.ErrInvalidParams, err.Error())
		return
	}
	agent, err := h.agentService.GetByID(ctx, GetAgentID(ctx))
	if err != nil {
		Fail(ctx, err)
		return
	}
	var res response.AgentResponse
	copier.Copy(&res, agent)
	res.VersionName = "v1.0.0"
	Success(ctx, res)
}

func (h *AgentController) GetAgentBind(ctx *gin.Context) {
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

	agentConfig := config.AgentConfig{
		ServerAddr: cfg.Server.TCP.Addr,
		AppId:      appID,
		AppSecret:  app.Secret,
	}

	data, err := json.Marshal(agentConfig)
	if err != nil {
		Fail(ctx, err)
		return
	}

	c := utils.Base64Encode(data)

	sh := "curl -s http://%s/file/%s -o /tmp/noah-cli && chmod +x /tmp/noah-cli && /tmp/noah-cli run -c %s"
	Success(ctx, response.AgentBindResponse{
		MacBind: fmt.Sprintf(sh, cfg.Server.HTTP.Addr, "noah-mac", c),
	})
}

func (h *AgentController) DeleteAgent(ctx *gin.Context) {
	agent, err := h.agentService.Delete(ctx, GetAgentID(ctx))
	if err != nil {
		Fail(ctx, err)
		return
	}

	// 通知客户端退出程序
	if err := session.GetSessionManager().SendCommand(agent.SessionID, packet.Command_EXIT); err != nil {
		logger.Info("agent logout fail", "err", err)
	}
	Success(ctx, nil)
}

func (h *AgentController) GetAgentMetric(ctx *gin.Context) {
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

	list, err := h.agentService.GetAgentMetric(ctx, GetAgentID(ctx), startTime, endTime)
	if err != nil {
		Fail(ctx, err)
		return
	}
	var agentInfoList []response.AgentMetricResponse
	if err := copier.CopyWithOption(&agentInfoList, list, copier.Option{
		Converters: []copier.TypeConverter{{
			SrcType: "",
			DstType: []response.DiskUsage{},
			Fn: func(src interface{}) (interface{}, error) {
				var dst []response.DiskUsage
				err := json.Unmarshal([]byte(src.(string)), &dst)
				return dst, err
			},
		}},
	}); err != nil {
		Fail(ctx, err)
		return
	}
	Success(ctx, agentInfoList)
}

func (h *AgentController) OpenPty(ctx *gin.Context) {
	logger.Info("open pty")
	agent, err := h.agentService.GetByID(ctx, GetAgentID(ctx))
	if err != nil {
		Fail(ctx, err)
		return
	}
	src, err := session.GetSessionManager().OpenTunnel(agent.SessionID, packet.OpenTunnel_PTY, "")
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

func (h *AgentController) GenerateV2raySubscribeLink(ctx *gin.Context) {
	var req request.AgentGenerateV2raySubscribeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Fail(ctx, errcode.ErrInvalidParams)
		return
	}
	agents, err := h.agentService.GetByIDs(ctx, req.Ids)
	if err != nil {
		Fail(ctx, err)
		return
	}

	cfg := config.Get()
	links := []string{}

	host, port, _ := net.SplitHostPort(cfg.Server.V2ray.Addr)

	for _, agent := range agents {
		addr := fmt.Sprintf(`{"add":"%s","id":"%s","net":"tcp","port":"%s","ps":"%s","scy":"auto","type":"none","v":"2"}`,
			host, agent.SessionID, port, cfg.Server.V2ray.Addr)
		links = append(links, "vmess://"+base64.StdEncoding.EncodeToString([]byte(addr)))
	}

	content := strings.Join(links, "\n")

	filename := fmt.Sprintf("v2ray-sub-%s.txt", utils.MD5(content))
	res := fmt.Sprintf("http://%s/file/%s", cfg.Server.HTTP.Addr, filename)
	if utils.FileExists("./temp/" + filename) {
		Success(ctx, res)
		return
	}

	file, err := os.Create("./temp/" + filename)
	if err != nil {
		Fail(ctx, err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		Fail(ctx, err)
		return
	}

	Success(ctx, res)
}
