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
	agentService service.AgentService
}

func NewPingHandler() conn.MessageHandler {
	return &pingHandler{
		agentService: service.GetAgentService(),
	}
}

func (p *pingHandler) Handle(ctx conn.Context) error {
	var ping packet.Ping
	if err := ctx.Unmarshal(&ping); err != nil {
		return err
	}
	logger.Info("receive ping msg", "data", ping.String())
	var agentStat model.AgentMetric
	copier.CopyWithOption(&agentStat, &ping, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	du, _ := json.Marshal(ping.DiskUsage)
	agentStat.DiskUsage = string(du)
	p.agentService.SaveAgentMetric(ctx, getSessionID(ctx), &agentStat)
	return nil
}

func (p *pingHandler) MessageType() packet.MessageType {
	return packet.MessageType_Ping
}
