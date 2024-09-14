package service

import (
	"noah/internal/server/logic/channel"
	"noah/internal/server/logic/client"
)

var (
	channelService        IChannelService
	clientServiceInstance IClientService
)

func LoadService() {
	clientServiceInstance = client.NewClientService()
	channelService = channel.NewChannelService()
}

func GetChannelService() IChannelService {
	return channelService
}

func GetClientService() IClientService {
	return clientServiceInstance
}
