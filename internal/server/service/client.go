package service

import (
	"noah/internal/server/dao"
	"noah/internal/server/request"
	"noah/internal/server/response"
	"time"
)

type IClientService interface {
	Save(client dao.Client) (id uint, err error)
	UpdateStatus(id uint, status int8)
	GetClient(id uint) (response.GetClientRes, error)
	GetClientPage(query request.ListClientQueryReq) (total int64, clients []response.ListClientRes)
	ScheduleUpdateStatus() error
	Delete(id uint) error
	Generate(serverAddr string, port string, osType int8, token string, filename string) (string, error)
	SaveSystemInfo(id uint, systemInfo request.CreateSystemInfoReq) error
	GetSystemInfo(id uint, start time.Time, end time.Time) ([]response.GetSystemInfoRes, error)
	CleanSystemInfo() error
	Count() (online int64, offline int64)
}
