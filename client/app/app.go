package app

import (
	"noah/client/app/environment"
	"noah/client/app/gateway"
	"noah/client/app/handler"
	"noah/client/app/service"
	"sync"

	"github.com/samber/do/v2"
)

type Client struct {
	Env      *environment.Environment
	Injector do.Injector
}

func NewClient(env *environment.Environment) *Client {
	i := Inject()

	do.ProvideValue(i, env)
	return &Client{
		Env:      env,
		Injector: i,
	}
}

func Inject() do.Injector {
	i := do.New()
	//init service
	service.Init(i)

	//gateway
	do.Provide(i, gateway.NewGateway)

	return i
}

func (c *Client) Start() {
	handler, err := handler.NewHandler(c.Injector)
	if err != nil {
		panic(err)
	}

	w := sync.WaitGroup{}
	w.Add(1)

	go func() {
		defer w.Done()
		handler.Connect()
	}()

	w.Wait()

}
