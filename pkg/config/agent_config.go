package config

type AgentConfig struct {
	ServerAddr string `json:"server_addr"` // 服务器地址
	AppId      uint64 `json:"app_id"`      // 应用ID
	AppSecret  string `json:"app_secret"`  // 应用密钥
	Version    uint32 `json:"version"`     // 应用版本

	DailTimeout       int `json:"dail_timeout"`       // 连接超时，单位秒
	ReconnectInterval int `json:"reconnect_interval"` // 重连间隔，单位秒
	HeartbeatInterval int `json:"heartbeat_interval"` // 心跳间隔，单位秒
}
