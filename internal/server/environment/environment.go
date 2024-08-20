package environment

import (
	"gopkg.in/yaml.v3"
	"os"
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

	return env, nil
}
