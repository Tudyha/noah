package vo

type DeviceFileRenamePostVo struct {
	Filename string `json:"filename" binding:"required"`
	Path     string `json:"path" binding:"required"`
}

type DeviceFileDeletePostVo struct {
	Path string `json:"path" binding:"required"`
	Type int8   `json:"type" binding:"required"`
}

type DeviceFileContentPostVo struct {
	Content string `json:"content" binding:"required"`
	Path    string `json:"path" binding:"required"`
}

type DeviceFileNewDirPostVo struct {
	Path string `json:"path" binding:"required"`
}
