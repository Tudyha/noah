package service

import (
	"noah/pkg/mux/message"
	"noah/pkg/request"
)

type IInformationService interface {
	LoadClientSpecs() (*request.CreateClientReq, error)
	GetSystemInfo() (*request.CreateClientStatReq, error)
	GetProcessList() ([]message.Process, error)
	KillProcess(pid int32) error
	GetNetworkInfo() ([]message.NetworkInfo, error)
	GetDockerContainerList() (res []message.DockerContainer, err error)
}
