package entitie

import "time"

type FileExplorer struct {
	Path     string    `json:"path"`
	Filename string    `json:"filename"`
	ModTime  time.Time `json:"modTime"`
	Type     uint8     `json:"type"`
	Size     int64     `json:"size"`
	Mod      string    `json:"mod"`
}

type FileExplorerQuery struct {
	Path        string `json:"path"`
	Op          string `json:"op"`
	Filename    string `json:"filename"`
	FileContent string `json:"fileContent"`
}
