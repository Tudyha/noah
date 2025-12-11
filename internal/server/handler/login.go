package handler

import (
	"noah/pkg/conn"
	"noah/pkg/logger"
	"noah/pkg/packet"

	"google.golang.org/protobuf/proto"
)

type loginHandler struct {
}

func NewLoginHandler() conn.MessageHandler {
	return &loginHandler{}
}

func (r *loginHandler) Handle(ctx conn.Context, msg proto.Message) error {
	logger.Info("login", "msg", msg)
	return ctx.Close()
}

func (r *loginHandler) MessageType() packet.MessageType {
	return packet.MessageType_Login
}

func (r *loginHandler) MessageBody() proto.Message {
	return &packet.Login{}
}
