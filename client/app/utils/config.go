package utils

import (
	"bytes"
	"encoding/json"
)

type Config struct {
	Version       string `json:"version"`
	VersionCode   int    `json:"version_code"`
	ServerPort    string `json:"server_port"`
	ServerAddress string `json:"server_address"`
	Token         string `json:"token"`
}

func ReadConfigFile(configFile []byte) *Config {
	configFile = bytes.NewBufferString(string(configFile)).Bytes()

	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		panic(err)
	}
	return &config
}
