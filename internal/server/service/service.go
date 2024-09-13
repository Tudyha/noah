package service

import (
	"noah/internal/server/logic/channel"
	"noah/internal/server/logic/client"
)

func LoadService() {
	clientServiceInstance = client.NewClientService()
	channelService = channel.NewChannelService()
}
