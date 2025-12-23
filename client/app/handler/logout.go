package handler

import (
	"log"
	"noah/pkg/conn"
	"noah/pkg/packet"
	"os"
)

type LogoutHandler struct {
}

func NewLogoutHandler() conn.MessageHandler {
	return &LogoutHandler{}
}

func (s *LogoutHandler) Handle(ctx conn.Context) error {
	log.Println("收到服务端退出消息，退出程序")
	os.Exit(0)
	return nil
}

func (s *LogoutHandler) MessageType() packet.MessageType {
	return packet.MessageType_Logout
}
