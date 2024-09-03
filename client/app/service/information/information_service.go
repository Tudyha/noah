package information

import (
	"os"
	"os/user"
	"runtime"
	"time"

	"noah/client/app/entitie"
	"noah/client/app/service"

	"noah/client/app/utils/network"
)

type Service struct {
}

func NewService() service.Information {
	return &Service{}
}

func (i Service) LoadClientSpecs() (*entitie.Client, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	username, err := user.Current()
	if err != nil {
		return nil, err
	}
	macAddress, err := network.GetMacAddress()
	if err != nil {
		return nil, err
	}
	return &entitie.Client{
		Hostname:    hostname,
		Username:    username.Name,
		UserID:      username.Username,
		OSName:      runtime.GOOS,
		OSArch:      runtime.GOARCH,
		MacAddress:  macAddress,
		IPAddress:   network.GetLocalIP(),
		Port:        "",
		FetchedUnix: time.Now().UTC().Unix(),
	}, nil
}
