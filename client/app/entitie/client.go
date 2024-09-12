package entitie

type Client struct {
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

type SystemInfo struct {
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

type Process struct {
	Pid        int32   `json:"pid"`
	Name       string  `json:"name"`
	Uids       []int32 `json:"uids"`
	Username   string  `json:"username"`
	Gids       []int32 `json:"gids"`
	Cpu        float64 `json:"cpu"`
	Memory     uint64  `json:"memory"`
	Command    string  `json:"command"`
	CreateTime int64   `json:"createTime"`
}
