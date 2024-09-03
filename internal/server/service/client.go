package service

import (
	"github.com/gorilla/websocket"
	"noah/internal/server/dto"
	"noah/internal/server/vo"
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
	Save(body dto.ClientPostDto) (id uint, err error)
	UpdateStatus(id uint, status int8)
	GetClient(query vo.ClientListQueryReq) (total int64, ClientDtos []dto.ClientDto)
	ScheduleUpdateStatus() error
	Delete(id uint) error
	AddConnection(id uint, connection *websocket.Conn) error
	SendCommand(id uint, commandStr string, parameter string) (result string, err error)
	Generate(serverAddr string, port string, osType int8, filename string) (string, error)
	Exit(id uint) error
}
