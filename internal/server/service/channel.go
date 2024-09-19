package service

import (
	"noah/internal/server/enum"
	"noah/internal/server/request"
	"noah/internal/server/response"

	"github.com/gorilla/websocket"
)

type IChannelService interface {
	NewClientWebsocketConn(id uint, connection *websocket.Conn) error
	Exit(id uint) error
	SendCommand(id uint, messageType enum.MessageType, data any) (string, error)
	NewChannel(id uint, channelReq request.CreateChannelReq, conn *websocket.Conn) error
	GetChannelList(clientId uint) ([]response.GetChannelListRes, error)
	DeleteChannel(id uint) (err error)
}
