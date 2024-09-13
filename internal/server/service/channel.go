package service

import (
	"github.com/gorilla/websocket"
	"noah/internal/server/enum"
	"noah/internal/server/logic/channel"
)

var (
	channelService IChannelService
)

func GetChannelService() IChannelService {
	return channelService
}

type IChannelService interface {
	NewChannel(channelType enum.ChannelType) (channel *channel.Channel)
	ClientConnect(channelId string, conn *websocket.Conn) error
}
