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

type loginHandler struct {
	clientService service.ClientService
	pub           *gochannel.GoChannel
}

func NewLoginHandler() conn.MessageHandler {
	return &loginHandler{
		clientService: service.GetClientService(),
		pub:           mq.GetPubSub(),
	}
}

func (h *loginHandler) Handle(ctx conn.Context) (err error) {
	defer func() {
		if err != nil {
			logger.Info("login fail, close conn", "err", err)
			ctx.GetConn().Close()
		}
	}()

	var loginReq packet.Login
	if err := ctx.Unmarshal(&loginReq); err != nil {
		return err
	}
	if err := h.clientService.VerifySign(ctx, loginReq.AppId, loginReq.Sign); err != nil {
		logger.Info("校验签名失败，断开连接", "appID", loginReq.AppId, "sign", loginReq.Sign)
		return err
	}

	var client model.Client

	if err := copier.Copy(&client, loginReq.ClientInfo); err != nil {
		logger.Info("复制数据失败", "err", err)
		return err
	}
	client.DeviceID = loginReq.DeviceId
	client.AppID = loginReq.AppId
	client.OsType = enum.ClientOsNameToOsTypeMap[client.OsName]

	remoteAddr := ctx.GetConn().RemoteAddr()
	remoteIP, port, _ := net.SplitHostPort(remoteAddr.String())
	client.RemoteIP = remoteIP
	client.Port = port
	client.RemoteIpCountry = ip.GetIPCountry(client.RemoteIP)
	client.SessionID = getSessionID(ctx)

	if err := h.clientService.Connect(ctx, &client); err != nil {
		logger.Info("创建客户端失败", "err", err)
		return err
	}

	h.pub.Publish(constant.MQ_TOPIC_CLIENT_ONLINE, message.NewMessage(watermill.NewUUID(), message.Payload(client.SessionID)))

	return nil
}

func (h *loginHandler) MessageType() packet.MessageType {
	return packet.MessageType_Login
}
