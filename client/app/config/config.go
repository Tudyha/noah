package config

type Config struct {
	ServerAddr        string `json:"server_addr"`        // 服务器地址
	AppId             uint64 `json:"app_id"`             // 应用ID
	AppSecret         string `json:"app_secret"`         // 应用密钥
	ReconnectInterval int    `json:"reconnect_interval"` // 重连间隔，单位秒
}
