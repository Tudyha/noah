package utils

import (
	"bytes"
	"encoding/json"
)

type Config struct {
	ServerPort    string `json:"server_port"`
	ServerAddress string `json:"server_address"`
	Token         string `json:"token"`
}

func ReadConfigFile(configFile []byte) *Config {
	// decoded, err := encode.DecodeBase64(bytes.NewBuffer(configFile).String())
	// if err != nil {
	// 	log.Fatal("error reading config file: ", err)
	// }

	configFile = bytes.NewBufferString(string(configFile)).Bytes()

	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		panic(err)
	}
	return &config
}
