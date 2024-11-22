package model

import "gorm.io/gorm"

type Tunnel struct {
	gorm.Model
	ClientId   uint   // 客户端id
	TunnelType uint8  // 通道类型
	ServerPort int    // 服务端端口
	ClientIp   string // 客户端ip
	ClientPort int    // 客户端端口
	Status     uint8  // 服务端状态
	FailReason string // 失败原因
}

func (Tunnel) TableName() string {
	return "tunnel"
}
