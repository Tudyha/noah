package handler

import (
	"net"
	"noah/internal/model"
	"noah/internal/service"
	"noah/pkg/conn"
	"noah/pkg/enum"
	"noah/pkg/logger"
	"noah/pkg/packet"

	"noah/pkg/ip"

	"github.com/jinzhu/copier"
)

type loginHandler struct {
	clientService service.ClientService
}

func NewLoginHandler() conn.MessageHandler {
	return &loginHandler{
		clientService: service.GetClientService(),
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
	if err := ctx.ShouldBindProto(&loginReq); err != nil {
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

	if err := h.clientService.Connect(ctx, &client); err != nil {
		logger.Info("创建客户端失败", "err", err)
		return err
	}

	return nil
}

func (h *loginHandler) MessageType() packet.MessageType {
	return packet.MessageType_Login
}
