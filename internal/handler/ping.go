package handler

import (
	"encoding/json"
	"noah/internal/model"
	"noah/internal/service"
	"noah/pkg/conn"
	"noah/pkg/logger"
	"noah/pkg/packet"

	"github.com/jinzhu/copier"
)

type pingHandler struct {
	clientService service.ClientService
}

func NewPingHandler() conn.MessageHandler {
	return &pingHandler{
		clientService: service.GetClientService(),
	}
}

func (p *pingHandler) Handle(ctx conn.Context) error {
	var ping packet.Ping
	if err := ctx.ShouldBindProto(&ping); err != nil {
		return err
	}
	logger.Info("receive ping msg", "data", ping.String())
	var clientStat model.ClientStat
	copier.CopyWithOption(&clientStat, &ping, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	du, _ := json.Marshal(ping.DiskUsage)
	clientStat.DiskUsage = string(du)
	// clientStat.ClientId = ctx.GetConn().GetClientID()
	p.clientService.SaveClientStat(ctx, &clientStat)
	return nil
}

func (p *pingHandler) MessageType() packet.MessageType {
	return packet.MessageType_Ping
}
