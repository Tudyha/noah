package config

type Config struct {
	ServerAddr string `json:"server_addr"`
	AppId      int32  `json:"app_id"`
	AppSecret  string `json:"app_secret"`
}
