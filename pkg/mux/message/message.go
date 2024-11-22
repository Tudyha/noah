package message

import "time"

type LinkInfo struct {
	Network string `json:"network"`
	Addr    string `json:"addr"`
}

type CmdInfo struct {
	Cmd  string `json:"cmd"`
	Data []byte `json:"data"`
}

type DownloadReq struct {
	Filename string `json:"filename"`
	Path     string `json:"path"`
}

type FileExplorer struct {
	Path     string    `json:"path"`
	Filename string    `json:"filename"`
	ModTime  time.Time `json:"modTime"`
	Type     uint8     `json:"type"`
	Size     int64     `json:"size"`
	Mod      string    `json:"mod"`
}

type FileExplorerReq struct {
	Path        string `json:"path"`
	Op          string `json:"op"`
	Filename    string `json:"filename"`
	FileContent string `json:"fileContent"`
}

type SystemInfoReq struct {
	SystemInfoType string `json:"systemInfoType"`
	Action         string `json:"action"`
	Params         string `json:"params"`
}

type SystemInfo struct {
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

type NetworkInfo struct {
	Fd     uint32          `json:"fd"`
	Family uint32          `json:"family"`
	Type   uint32          `json:"type"`
	Laddr  NetworkInfoAddr `json:"localaddr"`
	Raddr  NetworkInfoAddr `json:"remoteaddr"`
	Status string          `json:"status"`
	Uids   []int32         `json:"uids"`
	Pid    int32           `json:"pid"`
}

type NetworkInfoAddr struct {
	IP   string `json:"ip"`
	Port uint32 `json:"port"`
}

type DockerContainer struct {
	ID         string `json:"Id"`
	Names      []string
	Image      string
	ImageID    string
	Command    string
	Created    int64
	Ports      []DockerContainerPort
	SizeRw     int64 `json:",omitempty"`
	SizeRootFs int64 `json:",omitempty"`
	Labels     map[string]string
	State      string
	Status     string
	HostConfig struct {
		NetworkMode string            `json:",omitempty"`
		Annotations map[string]string `json:",omitempty"`
	}
}

type DockerContainerPort struct {
	IP          string `json:"IP,omitempty"`
	PrivatePort uint16 `json:"PrivatePort"`
	PublicPort  uint16 `json:"PublicPort,omitempty"`
	Type        string `json:"Type"`
}
