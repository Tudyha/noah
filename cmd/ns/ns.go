package main

import (
	"noah/internal/server"
	"os"

	"noah/internal/server/environment"

	"gopkg.in/yaml.v3"
)

func main() {
	env, err := LoadEnvironment()
	if err != nil {
		panic(err)
	}

	// new server
	s := server.NewServer(&env)

	//start server
	s.Run()
}

func LoadEnvironment() (environment.Environment, error) {
	var env environment.Environment

	file, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		return env, err
	}

	err = yaml.Unmarshal(file, &env)
	if err != nil {
		return env, err
	}

	//如果环境变量设置了配置，则覆盖配置文件
	if os.Getenv("HOST") != "" {
		env.Server.Host = os.Getenv("HOST")
	}
	// if os.Getenv("ADMIN_PASSWORD") != "" {
	// 	env.Admin.Password = os.Getenv("ADMIN_PASSWORD")
	// }

	return env, nil
}
