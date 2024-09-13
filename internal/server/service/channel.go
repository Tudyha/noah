package service

import (
	"noah/internal/server/enum"
	"noah/internal/server/logic/channel"

	"github.com/gorilla/websocket"
)

var (
	channelService IChannelService
)

func GetChannelService() IChannelService {
	return channelService
}

type IChannelService interface {
	NewChannel(channelType enum.ChannelType, serverPort string) (channel *channel.Channel)
	ClientConnect(channelId string, conn *websocket.Conn) error
}
