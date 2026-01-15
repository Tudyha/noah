package handler

import (
	"net"
	"noah/internal/model"
	"noah/internal/mq"
	"noah/internal/service"
	"noah/pkg/conn"
	"noah/pkg/constant"
	"noah/pkg/enum"
	"noah/pkg/logger"
	"noah/pkg/packet"

	"noah/pkg/ip"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/jinzhu/copier"
)

type authHandler struct {
	agentService service.AgentService
	pub          *gochannel.GoChannel
}

func NewAuthHandler() conn.MessageHandler {
	return &authHandler{
		agentService: service.GetAgentService(),
		pub:          mq.GetPubSub(),
	}
}

func (h *authHandler) Handle(ctx conn.Context) (err error) {
	defer func() {
		if err != nil {
			logger.Info("login fail, close conn", "err", err)
			ctx.GetConn().Close()
		}
	}()

	var loginReq packet.Auth
	if err := ctx.Unmarshal(&loginReq); err != nil {
		return err
	}
	if err := h.agentService.VerifySign(ctx, loginReq.AppId, loginReq.Sign); err != nil {
		logger.Info("校验签名失败，断开连接", "appID", loginReq.AppId, "sign", loginReq.Sign)
		return err
	}

	var agent model.Agent

	if err := copier.Copy(&agent, loginReq.AgentInfo); err != nil {
		logger.Info("复制数据失败", "err", err)
		return err
	}
	agent.DeviceID = loginReq.DeviceId
	agent.AppID = loginReq.AppId
	agent.OsType = enum.AgentOsNameToOsTypeMap[agent.OsName]
	agent.Version = uint32(loginReq.Version)

	remoteAddr := ctx.GetConn().RemoteAddr()
	remoteIP, port, _ := net.SplitHostPort(remoteAddr.String())
	agent.RemoteIP = remoteIP
	agent.Port = port
	agent.RemoteIpCountry = ip.GetIPCountry(agent.RemoteIP)
	agent.SessionID = getSessionID(ctx)

	if err := h.agentService.Connect(ctx, &agent); err != nil {
		logger.Info("创建客户端失败", "err", err)
		return err
	}

	h.pub.Publish(constant.MQ_TOPIC_CLIENT_ONLINE, message.NewMessage(watermill.NewUUID(), message.Payload(agent.SessionID)))

	return nil
}

func (h *authHandler) MessageType() packet.MessageType {
	return packet.MessageType_Login
}
