package main

import (
	_ "embed"
	"noah/client/app"
	"noah/client/app/environment"
	"noah/client/app/utils"
)

var (
	Version = "dev"
)

//go:embed config.json
var configFile []byte

func main() {
	config := utils.ReadConfigFile(configFile)

	// ui.ShowMenu(Version, config.ServerAddress, config.Port)

	app.New(environment.Load(config.ServerAddress, config.ServerPort, config.Token)).Run()
}
