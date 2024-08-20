package dto

import "time"

type DevicePostDto struct {
	Hostname   string `json:"hostname" binding:"required"`
	Username   string `json:"username" binding:"required"`
	UserID     string `json:"user_id" binding:"required"`
	OSName     string `json:"os_name" binding:"required"`
	OSArch     string `json:"os_arch" binding:"required"`
	MacAddress string `json:"mac_address" binding:"required"`
	IPAddress  string `json:"ip_address"`
	Port       string `json:"port"`
}

type DeviceDto struct {
	ID             uint      `json:"id" binding:"required"`
	Hostname       string    `json:"hostname" binding:"required"`
	Username       string    `json:"username" binding:"required"`
	UserID         string    `json:"user_id" binding:"required"`
	OSName         string    `json:"os_name" binding:"required"`
	OSArch         string    `json:"os_arch" binding:"required"`
	MacAddress     string    `json:"mac_address" binding:"required"`
	IPAddress      string    `json:"ip_address"`
	Port           string    `json:"port"`
	Status         int8      `json:"status" comment:"0-offline 1-online"`
	LastOnlineTime time.Time `json:"last_online_time" comment:"最后上线时间"`
}

type DeviceListQueryDto struct {
	PageQuery
	Hostname string `form:"hostname"`
}
