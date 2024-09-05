package service

import (
	"github.com/gorilla/websocket"
	"noah/internal/server/dto"
	"noah/internal/server/vo"
	"time"
)

var (
	clientServiceInstance IClientService
)

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
	Generate(serverAddr string, port string, osType int8, token string, filename string) (string, error)
	Exit(id uint) error
	SaveSystemInfo(id uint, systemInfo dto.SystemInfoReq) error
	GetSystemInfo(id uint, start time.Time, end time.Time) ([]dto.SystemInfoRes, error)
	CleanSystemInfo() error
}
