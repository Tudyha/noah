package request

type CreateTunnelReq struct {
	TunnelType uint8  `json:"tunnelType" binding:"required"`
	ServerPort int    `json:"serverPort" binding:"required"`
	TargetAddr string `json:"targetAddr"`
}
