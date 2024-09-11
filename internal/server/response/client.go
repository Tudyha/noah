package response

import (
	"noah/internal/server/enum"
	"time"
)

type ListClientRes struct {
	ID             uint        `json:"id" binding:"required"`
	Hostname       string      `json:"hostname" binding:"required"`
	Username       string      `json:"username" binding:"required"`
	Gid            string      `json:"gId" binding:"required"`
	Uid            string      `json:"userId" binding:"required"`
	OsType         enum.OSType `json:"osType"`
	OSName         string      `json:"osName" binding:"required"`
	OSArch         string      `json:"osArch" binding:"required"`
	MacAddress     string      `json:"macAddress" binding:"required"`
	CpuCores       int32       `json:"cpuCores"`
	CpuModelName   string      `json:"cpuModelName"`
	CpuFamily      string      `json:"cpuFamily"`
	MemoryTotal    float64     `json:"memoryTotal"`
	DiskTotal      float64     `json:"diskTotal"`
	RemoteIp       string      `json:"remoteIp"`
	LocalIp        string      `json:"localIp"`
	Status         int8        `json:"status" comment:"0-offline 1-online"`
	LastOnlineTime time.Time   `json:"lastOnlineTime" comment:"最后上线时间"`
}

type GetClientRes struct {
	ID             uint        `json:"id" binding:"required"`
	Hostname       string      `json:"hostname" binding:"required"`
	Username       string      `json:"username" binding:"required"`
	Gid            string      `json:"gId" binding:"required"`
	Uid            string      `json:"userId" binding:"required"`
	OsType         enum.OSType `json:"osType"`
	OSName         string      `json:"osName" binding:"required"`
	OSArch         string      `json:"osArch" binding:"required"`
	MacAddress     string      `json:"macAddress" binding:"required"`
	CpuCores       int32       `json:"cpuCores"`
	CpuModelName   string      `json:"cpuModelName"`
	CpuFamily      string      `json:"cpuFamily"`
	MemoryTotal    float64     `json:"memoryTotal"`
	DiskTotal      float64     `json:"diskTotal"`
	RemoteIp       string      `json:"remoteIp"`
	LocalIp        string      `json:"localIp"`
	Status         int8        `json:"status" comment:"0-offline 1-online"`
	LastOnlineTime time.Time   `json:"lastOnlineTime" comment:"最后上线时间"`
}
