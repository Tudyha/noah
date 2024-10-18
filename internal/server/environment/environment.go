package environment

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Environment struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Admin    AdminConfig    `yaml:"admin"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type AdminConfig struct {
	Password string `yaml:"password"`
}

func LoadEnvironment() (Environment, error) {
	var env Environment

	file, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		return env, err
	}

	err = yaml.Unmarshal(file, &env)
	if err != nil {
		return env, err
	}

	//如果环境变量设置了配置，则覆盖配置文件
	if os.Getenv("SERVER_PORT") != "" {
		env.Server.Port = os.Getenv("SERVER_PORT")
	}
	if os.Getenv("ADMIN_PASSWORD") != "" {
		env.Admin.Password = os.Getenv("ADMIN_PASSWORD")
	}

	return env, nil
}
