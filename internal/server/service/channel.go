package service

import (
	"noah/internal/server/enum"
	"noah/internal/server/logic/channel"

	"github.com/gorilla/websocket"
)

type IChannelService interface {
	NewClientWebsocketConn(id uint, connection *websocket.Conn) error
	Exit(id uint) error
	SendCommand(id uint, messageType enum.MessageType, data any) (string, error)
	NewChannel(id uint, channelType enum.ChannelType, conn *websocket.Conn, serverPort string) (channel *channel.Channel, err error)
}
