package main

import (
	"noah/client/app"
	"noah/client/app/environment"
)

func main() {
	env := environment.Environment{
		Server: environment.ServerConfig{
			Host: "http://127.0.0.1:8080",
		}}
	c := app.NewClient(&env)
	c.Start()
}
