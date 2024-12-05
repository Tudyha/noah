package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"noah/client/app"
	"noah/client/app/environment"
)

//go:embed config.json
var configFile []byte

func main() {
	var env environment.Environment
	if err := json.Unmarshal(configFile, &env); err != nil {
		log.Fatalln(err)
	}
	cli := app.NewClient(&env)
	cli.Start()
}
