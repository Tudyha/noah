package request

type CreateTunnelReq struct {
	TunnelType uint8  `json:"tunnelType" binding:"required"`
	ClientIp   string `json:"clientIp"`
	ClientPort int    `json:"clientPort"`
	ServerPort int    `json:"serverPort" binding:"required"`
}
