package service

import (
	"noah/pkg/request"
	"noah/pkg/response"
)

type ITunnelService interface {
	NewTunnel(id uint, TunnelReq request.CreateTunnelReq) error
	GetTunnelList(clientId uint) ([]response.GetTunnelListRes, error)
	DeleteTunnel(id uint) (err error)
	DeleteTunnelByClientId(clientId uint) (err error)
}
