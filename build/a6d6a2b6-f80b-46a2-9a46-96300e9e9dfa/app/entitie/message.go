package entitie

type CommandReq struct {
	Command string `json:"command,omitempty"`
}

type DownloadReq struct {
	Filename string `json:"filename"`
	Path     string `json:"path"`
}

type SystemInfoReq struct {
	SystemInfoType string `json:"systemInfoType"`
	Action         string `json:"action"`
	Params         string `json:"params"`
}
