package service

import (
	"github.com/samber/do/v2"
	"noah/internal/server/logic/channel"
	"noah/internal/server/logic/client"
)

var (
	injector do.Injector
)

func LoadService(i do.Injector) {
	injector = i
	do.Provide(i, client.NewClientService)
	do.Provide(i, channel.NewChannelService)
}

func GetChannelService() IChannelService {
	return do.MustInvoke[*channel.Service](injector)
}

func GetClientService() IClientService {
	return do.MustInvoke[*client.Service](injector)
}
