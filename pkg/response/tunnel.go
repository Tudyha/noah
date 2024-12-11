package response

type GetTunnelListRes struct {
	ID         uint   `json:"id"`
	TunnelType uint8  `json:"tunnelType" binding:"required"`
	ClientIp   string `json:"clientIp"`
	ClientPort int    `json:"clientPort"`
	ServerPort int    `json:"serverPort"`
	Status     uint8  `json:"status"`
	FailReason string `json:"failReason"`
	TargetAddr string `json:"targetAddr"`
}
