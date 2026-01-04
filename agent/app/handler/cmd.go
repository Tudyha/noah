package handler

import (
	"log"
	"noah/pkg/conn"
	"noah/pkg/packet"
	"os"
)

type CommandHandler struct {
}

func NewCommandHandler() conn.MessageHandler {
	return &CommandHandler{}
}

func (c *CommandHandler) Handle(ctx conn.Context) error {
	var command packet.Command
	if err := ctx.Unmarshal(&command); err != nil {
		return err
	}
	switch command.Cmd {
	case packet.Command_EXIT:
		log.Println("收到服务端退出程序消息")
		os.Exit(0)
	}
	return nil
}

func (c *CommandHandler) MessageType() packet.MessageType {
	return packet.MessageType_Command
}
