package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	String() string
}

type App struct {
	servers []Server
}

func NewApp(servers ...Server) *App {
	return &App{
		servers: servers,
	}
}

func (a *App) Run() {
	ctx := context.Background()
	for _, server := range a.servers {
		go func() {
			if err := server.Start(ctx); err != nil {
				fmt.Println("start server", server.String(), "failed:", err)
				os.Exit(1)
			}
		}()
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func (a *App) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, server := range a.servers {
		if err := server.Stop(ctx); err != nil {
			log.Println("stop server", server.String(), "failed:", err)
		}
	}
}
