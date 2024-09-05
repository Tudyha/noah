package service

import "github.com/gorilla/websocket"

var (
	ptyServiceInstance IPtyService
)

func GetPtyService() IPtyService {
	return ptyServiceInstance
}

type IPtyService interface {
	NewPtyChannel(channelId string, conn *websocket.Conn) error
	NewPtyClient(channelId string, connection *websocket.Conn) error
	PtyClientRead(channelId string) (messageType int, data []byte, err error)
	PtyClientWrite(channelId string, messageType int, data []byte) error
	ClosePtyClient(channelId string) error
}
