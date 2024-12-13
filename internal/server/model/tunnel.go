package model

import "gorm.io/gorm"

type Tunnel struct {
	gorm.Model
	ClientId   uint   // 客户端id
	TunnelType uint8  // 通道类型
	ServerPort int    // 服务端端口
	TargetAddr string // 目标地址
	Status     uint8  // 服务端状态
	FailReason string // 失败原因
	Cipher     string // 加密方式
	Password   string // 密码
}

func (Tunnel) TableName() string {
	return "tunnel"
}
