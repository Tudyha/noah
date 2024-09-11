package dto

import (
	"time"
)

type ClientPostDto struct {
	Hostname   string `json:"hostname" binding:"required"`
	Username   string `json:"username" binding:"required"`
	UserID     string `json:"userId" binding:"required"`
	OSName     string `json:"osName" binding:"required"`
	OSArch     string `json:"osArch" binding:"required"`
	MacAddress string `json:"macAddress" binding:"required"`
	IPAddress  string `json:"ipAddress"`
	Port       string `json:"port"`
}

type Command struct {
	Command   string `json:"command,omitempty"`
	Parameter string `json:"parameter,omitempty"`
	Response  []byte `json:"response,omitempty"`
	HasError  bool   `json:"has_error,omitempty"`
}

type RespondCommandRequestBody struct {
	ClientID uint   `json:"client_id,omitempty"`
	Response []byte `json:"response"`
	HasError bool   `json:"has_error,omitempty"`
}

type SystemInfoReq struct {
	CpuUsage          float64 `json:"cpuUsage"`
	MemoryTotal       float64 `json:"memoryTotal"`
	MemoryUsed        float64 `json:"memoryUsed"`
	MemoryFree        float64 `json:"memoryFree"`
	MemoryUsedPercent float64 `json:"memoryPercent"`
	MemoryAvailable   float64 `json:"memoryAvailable"`
	DiskTotal         float64 `json:"diskTotal"`
	DiskFree          float64 `json:"diskFree"`
	DiskUsed          float64 `json:"diskUsed"`
	BandwidthIn       float64 `json:"bandwidthIn"`
	BandwidthOut      float64 `json:"bandwidthOut"`
}

type SystemInfoRes struct {
	CpuUsage          float64   `json:"cpuUsage"`
	MemoryTotal       float64   `json:"memoryTotal"`
	MemoryUsed        float64   `json:"memoryUsed"`
	MemoryFree        float64   `json:"memoryFree"`
	MemoryUsedPercent float64   `json:"memoryPercent"`
	MemoryAvailable   float64   `json:"memoryAvailable"`
	DiskTotal         float64   `json:"diskTotal"`
	DiskFree          float64   `json:"diskFree"`
	DiskUsed          float64   `json:"diskUsed"`
	BandwidthIn       float64   `json:"bandwidthIn"`
	BandwidthOut      float64   `json:"bandwidthOut"`
	CreatedAt         time.Time `json:"createdAt"`
}
