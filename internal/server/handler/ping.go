package handler

import (
	"noah/pkg/conn"
	"noah/pkg/logger"
	"noah/pkg/packet"
)

type pingHandler struct {
}

func NewPingHandler() conn.MessageHandler {
	return &pingHandler{}
}

func (p *pingHandler) Handle(ctx conn.Context) error {
	logger.Info("receive ping msg", "ConnID", ctx.GetConn().GetID())
	return nil
}

func (p *pingHandler) MessageType() packet.MessageType {
	return packet.MessageType_Ping
}
