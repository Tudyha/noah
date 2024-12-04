package service

import (
	"noah/internal/server/model"
	"noah/pkg/request"
	"noah/pkg/response"
	"time"
)

type IClientService interface {
	Save(client model.Client) (id uint, err error)
	// UpdateStatus(id uint, status int8)
	GetClient(id uint) (response.GetClientRes, error)
	GetClientPage(query request.ListClientQueryReq) (total int64, clients []response.ListClientRes)
	Delete(id uint) error
	SaveClientStat(id uint, systemInfo request.CreateClientStatReq) error
	GetClientStat(id uint, start time.Time, end time.Time) ([]response.GetClientStatRes, error)
	Count() (online int64, offline int64)
}
