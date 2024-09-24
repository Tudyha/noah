package service

import (
	"noah/internal/server/gateway"
	"noah/internal/server/logic/channel"
	"noah/internal/server/logic/client"
)

var (
	channelService        IChannelService
	clientServiceInstance IClientService
)

func LoadService(gateway *gateway.Gateway) {
	clientServiceInstance = client.NewClientService()
	channelService = channel.NewChannelService(gateway)
}

func GetChannelService() IChannelService {
	return channelService
}

func GetClientService() IClientService {
	return clientServiceInstance
}
