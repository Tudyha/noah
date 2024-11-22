package service

import (
	"noah/internal/server/service/logic/client"
	"noah/internal/server/service/logic/tunnel"
	"noah/internal/server/service/logic/user"

	"github.com/samber/do/v2"
)

func Init(i do.Injector) {
	do.Provide(i, func(i do.Injector) (IUserService, error) {
		return user.NewUserService(i)
	})

	do.Provide(i, func(i do.Injector) (IClientService, error) {
		return client.NewClientService(i)
	})

	do.Provide(i, func(i do.Injector) (ITunnelService, error) {
		return tunnel.NewTunnelService(i)
	})
}
