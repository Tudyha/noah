package request

type CreateClientReq struct {
	Hostname     string `json:"hostname"`
	Username     string `json:"username"`
	Gid          string `json:"gid"`
	Uid          string `json:"userId"`
	OSName       string `json:"osName"`
	OSArch       string `json:"osArch"`
	MacAddress   string `json:"macAddress"`
	IPAddress    string `json:"ipAddress"`
	Port         string `json:"port"`
	CpuCores     int32  `json:"cpuCores"`
	CpuModelName string `json:"cpuModelName"`
	CpuFamily    string `json:"cpuFamily"`
	MemoryTotal  uint64 `json:"memoryTotal"`
	DiskTotal    uint64 `json:"diskTotal"`
}

type ListClientQueryReq struct {
	PageQuery
	Hostname string `form:"hostname"`
	Status   int8   `form:"status"`
}

type CreateClientStatReq struct {
	CpuUsage          float64 `json:"cpuUsage"`
	MemoryTotal       uint64  `json:"memoryTotal"`
	MemoryUsed        uint64  `json:"memoryUsed"`
	MemoryFree        uint64  `json:"memoryFree"`
	MemoryUsedPercent float64 `json:"memoryPercent"`
	MemoryAvailable   uint64  `json:"memoryAvailable"`
	DiskTotal         uint64  `json:"diskTotal"`
	DiskFree          uint64  `json:"diskFree"`
	DiskUsed          uint64  `json:"diskUsed"`
	BandwidthIn       float64 `json:"bandwidthIn"`
	BandwidthOut      float64 `json:"bandwidthOut"`
}

// SendCommandReq 发送命令请求
type SendCommandReq struct {
	ID        uint   `json:"id" binding:"required" description:"唯一标识"`
	Command   string `json:"command" binding:"required" description:"命令"`
	Parameter string `json:"parameter,omitempty" description:"参数"`
}

// ClientGenerateReq 生成客户端请求
type ClientGenerateReq struct {
	ServerAddr string `json:"serverAddr" binding:"required" description:"服务器地址"`
	Port       string `json:"port" binding:"required" description:"端口号"`
	OsType     int8   `json:"osType,omitempty" description:"操作系统类型"`
	Filename   string `json:"filename,omitempty" description:"文件名"`
}

// ClientFileRenameReq 重命名文件请求
type ClientFileRenameReq struct {
	Filename string `json:"filename" binding:"required" description:"文件名"`
	Path     string `json:"path" binding:"required" description:"文件路径"`
}

// ClientFileDeleteReq 删除文件请求
type ClientFileDeleteReq struct {
	Path string `json:"path" binding:"required" description:"文件路径"`
}

// ClientFileContentReq 文件内容请求
type ClientFileContentReq struct {
	Content string `json:"content" binding:"required" description:"文件内容"`
	Path    string `json:"path" binding:"required" description:"文件路径"`
}

// ClientFileNewDirReq 创建新目录请求
type ClientFileNewDirReq struct {
	Path string `json:"path" binding:"required" description:"目录路径"`
}

type CreateSystemInfoReq struct {
	CpuUsage          float64 `json:"cpuUsage"`
	MemoryTotal       uint64  `json:"memoryTotal"`
	MemoryUsed        uint64  `json:"memoryUsed"`
	MemoryFree        uint64  `json:"memoryFree"`
	MemoryUsedPercent float64 `json:"memoryPercent"`
	MemoryAvailable   uint64  `json:"memoryAvailable"`
	DiskTotal         uint64  `json:"diskTotal"`
	DiskFree          uint64  `json:"diskFree"`
	DiskUsed          uint64  `json:"diskUsed"`
	BandwidthIn       float64 `json:"bandwidthIn"`
	BandwidthOut      float64 `json:"bandwidthOut"`
}
