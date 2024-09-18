package request

import "noah/internal/server/enum"

type CreateChannelReq struct {
	ChannelType enum.ChannelType `json:"channelType" binding:"required"`
	ClientIp    string           `json:"clientIp" binding:"required"`
	ClientPort  int              `json:"clientPort" binding:"required"`
	ServerPort  int              `json:"serverPort" binding:"required"`
}
