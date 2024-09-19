package request

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
	PageQuery
	Hostname string `form:"hostname"`
	Status   int8   `form:"status"`
}

type GetFileExplorerQueryReq struct {
	Path        string `json:"path"`
	Op          string `json:"op"`
	Filename    string `json:"filename"`
	FileContent string `json:"file_content"`
}
