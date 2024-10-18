package request

import "noah/internal/server/enum"

type CreateChannelReq struct {
	ChannelType enum.ChannelType `json:"channelType" binding:"required"`
	ClientIp    string           `json:"clientIp"`
	ClientPort  int              `json:"clientPort"`
	ServerPort  int              `json:"serverPort" binding:"required"`
}
