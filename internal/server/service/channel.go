package service

import (
	"noah/internal/server/enum"

	"github.com/gorilla/websocket"
)

type IChannelService interface {
	NewClientWebsocketConn(id uint, connection *websocket.Conn) error
	Exit(id uint) error
	SendCommand(id uint, messageType enum.MessageType, data any) (string, error)
	NewChannel(id uint, channelType enum.ChannelType, conn *websocket.Conn, serverPort string, clientAddr string) error
}
