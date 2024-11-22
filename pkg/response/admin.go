package response

type GetDashboardRes struct {
	ClientOnlineCount  int64   `json:"clientOnlineCount"`
	ClientOfflineCount int64   `json:"clientOfflineCount"`
	MemoryTotal        float64 `json:"memoryTotal"`
	MemoryUsed         float64 `json:"memoryUsed"`
	MemoryFree         float64 `json:"memoryFree"`
	MemoryUsedPercent  float64 `json:"memoryPercent"`
	MemoryAvailable    float64 `json:"memoryAvailable"`
	MemoryRemain       float64 `json:"memoryRemain"`
	DiskTotal          float64 `json:"diskTotal"`
	DiskFree           float64 `json:"diskFree"`
	DiskUsed           float64 `json:"diskUsed"`
	Hostname           string  `json:"hostname"`
	Uptime             uint64  `json:"uptime"`
	BootTime           uint64  `json:"bootTime"`
	OS                 string  `json:"os"`              // ex: freebsd, linux
	Platform           string  `json:"platform"`        // ex: ubuntu, linuxmint
	PlatformFamily     string  `json:"platformFamily"`  // ex: debian, rhel
	PlatformVersion    string  `json:"platformVersion"` // version of the complete OS
	KernelVersion      string  `json:"kernelVersion"`   // version of the OS kernel (if available)
	KernelArch         string  `json:"kernelArch"`      // native cpu architecture queried at runtime, as returned by `uname -m` or empty string in case of error
}
