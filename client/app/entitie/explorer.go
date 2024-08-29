package entitie

import "time"

type FileExplorer struct {
	Path     string    `json:"path"`
	Filename string    `json:"filename"`
	ModTime  time.Time `json:"mod_time"`
	Type     uint8     `json:"type"`
}

type FileExplorerQuery struct {
	Path        string `json:"path"`
	Op          string `json:"op"`
	Filename    string `json:"filename"`
	FileContent string `json:"file_content"`
}
