package service

import (
	"noah/internal/server/request"
	"noah/internal/server/response"

	"github.com/gorilla/websocket"
)

type IChannelService interface {
	NewChannel(id uint, channelReq request.CreateChannelReq, conn *websocket.Conn) error
	GetChannelList(clientId uint) ([]response.GetChannelListRes, error)
	DeleteChannel(id uint) (err error)
}
