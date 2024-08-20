package service

import (
	"github.com/gorilla/websocket"
)

var (
	clientServiceInstance IClientService
)

func RegisterClientService(i IClientService) {
	clientServiceInstance = i
}

func GetClientService() IClientService {
	return clientServiceInstance
}

type IClientService interface {
	AddConnection(id uint, connection *websocket.Conn) error
	// AddPtyConnection(channelId string, connection *websocket.Conn) error
	SendCommand(id uint, commandStr string, parameter string) (result string, err error)
	// PtyRead(channelId string) (messageType int, data []byte, err error)
	// PtyWrite(channelId string, messageType int, data []byte) error
	// ClosePtyConnection(channelId string) error
	Generate(serverAddr string, port string, osType int8, filename string) (string, error)
}
