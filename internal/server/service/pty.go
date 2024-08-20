package service

import "github.com/gorilla/websocket"

var (
	ptyServiceInstance IPtyService
)

func RegisterPtyService(i IPtyService) {
	ptyServiceInstance = i
}

func GetPtyService() IPtyService {
	return ptyServiceInstance
}

type IPtyService interface {
	NewPtyClient(channelId string, conn *websocket.Conn) error
	AddPtyConnection(channelId string, connection *websocket.Conn) error
	PtyRead(channelId string) (messageType int, data []byte, err error)
	PtyWrite(channelId string, messageType int, data []byte) error
	ClosePtyConnection(channelId string) error
}
