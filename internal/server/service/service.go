package service

import (
	"noah/internal/server/logic/client"
	"noah/internal/server/logic/pty"
)

func LoadService() {
	clientServiceInstance = client.NewClientService()
	ptyServiceInstance = pty.NewPtyService()
}
