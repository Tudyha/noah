package request

type CreateTunnelReq struct {
	TunnelType uint8  `json:"tunnelType" binding:"required"`
	ServerPort int    `json:"serverPort" binding:"required"`
	TargetAddr string `json:"targetAddr"`
	Cipher     string `json:"cipher"`
	Password   string `json:"password"`
}
