package app

import (
	"context"
	"log"
	"noah/client/app/environment"
	"noah/client/app/gateway/client"
	"noah/client/app/handler"
	"noah/client/app/service"
	"noah/client/app/service/information"
	"noah/client/app/service/pty"
	"noah/client/app/service/terminal"
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
		Information: information.NewService(),
		Terminal:    terminal.NewService(),
		Pty:         pty.NewService(),
	}

	return &App{
		handler.NewHandler(configuration, clientGateway, clientServices),
	}
}

func (a *App) Run() {
	id, err := a.Handler.SendDeviceSpecs()
	a.Handler.ClientID = id
	if err != nil {
		log.Fatal("error running client: ", err)
	}

	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		a.Handler.KeepConnection()
		return nil
	})

	g.Go(func() error {
		a.Handler.HandleCommand()
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Fatal("error running client: ", err)
	}
}
