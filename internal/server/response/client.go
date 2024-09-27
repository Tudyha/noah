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
	MemoryTotal    string      `json:"memoryTotal"`
	DiskTotal      string      `json:"diskTotal"`
	RemoteIp       string      `json:"remoteIp"`
	LocalIp        string      `json:"localIp"`
	Status         int8        `json:"status" comment:"0-offline 1-online"`
	LastOnlineTime time.Time   `json:"lastOnlineTime" comment:"最后上线时间"`
}

type GetFileExplorerRes struct {
	Path     string    `json:"path"`
	Filename string    `json:"filename"`
	ModTime  time.Time `json:"modTime"`
	Type     uint8     `json:"type"`
	Size     int64     `json:"size"`
	Mod      string    `json:"mod"`
}

type GetClientProcessRes struct {
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

type GetChannelListRes struct {
	ID          uint               `json:"id"`
	ChannelType enum.ChannelType   `json:"channelType"`
	ClientIp    string             `json:"clientIp"`
	ClientPort  int                `json:"clientPort"`
	ServerPort  int                `json:"serverPort"`
	Status      enum.ChannelStatus `json:"status"`
	FailReason  string             `json:"failReason"`
}

type GetClientNetworkInfoRes struct {
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

type GetClientDockerContainerRes struct {
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

type GetSystemInfoRes struct {
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
