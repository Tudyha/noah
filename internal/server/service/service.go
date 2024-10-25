package service

import (
	"noah/internal/server/logic/channel"
	"noah/internal/server/logic/client"

	"github.com/samber/do/v2"
)

var (
	injector do.Injector
)

func LoadService(i do.Injector) {
	injector = i
	clientService := client.NewClientService(i)
	do.ProvideValue(i, clientService)

	channelService := channel.NewChannelService(i)
	do.ProvideValue(i, channelService)
}

func GetChannelService() IChannelService {
	return do.MustInvoke[*channel.Service](injector)
}

func GetClientService() IClientService {
	return do.MustInvoke[*client.Service](injector)
}
