package request

import "noah/internal/server/dto"

type CreateClientReq struct {
	Hostname     string  `json:"hostname"`
	Username     string  `json:"username"`
	Gid          string  `json:"gid"`
	Uid          string  `json:"userId"`
	OSName       string  `json:"osName"`
	OSArch       string  `json:"osArch"`
	MacAddress   string  `json:"macAddress"`
	IPAddress    string  `json:"ipAddress"`
	Port         string  `json:"port"`
	CpuCores     int32   `json:"cpuCores"`
	CpuModelName string  `json:"cpuModelName"`
	CpuFamily    string  `json:"cpuFamily"`
	MemoryTotal  float64 `json:"memoryTotal"`
	DiskTotal    float64 `json:"diskTotal"`
}

type ListClientQueryReq struct {
	dto.PageQuery
	Hostname string `form:"hostname"`
	Status   int8   `form:"status"`
}
