package service

import (
	"github.com/gorilla/websocket"
	"noah/internal/server/dao"
	"noah/internal/server/dto"
	"noah/internal/server/request"
	"noah/internal/server/response"
	"time"
)

var (
	clientServiceInstance IClientService
)

func GetClientService() IClientService {
	return clientServiceInstance
}

type IClientService interface {
	Save(client dao.Client) (id uint, err error)
	UpdateStatus(id uint, status int8)
	GetClient(id uint) (response.GetClientRes, error)
	GetClientPage(query request.ListClientQueryReq) (total int64, clients []response.ListClientRes)
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
