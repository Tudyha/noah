package vo

import "noah/internal/server/dto"

// ClientPostReq 客户端注册请求
type ClientPostReq struct {
	Hostname    string `json:"hostname,omitempty" description:"主机名"`
	Username    string `json:"username,omitempty" description:"用户名"`
	UserID      string `json:"userId,omitempty" description:"用户ID"`
	OSName      string `json:"osName,omitempty" description:"操作系统名称"`
	OSArch      string `json:"osArch,omitempty" description:"操作系统架构"`
	MacAddress  string `json:"macAddress,omitempty" description:"MAC地址"`
	IPAddress   string `json:"ipAddress,omitempty" description:"IP地址"`
	Port        string `json:"port,omitempty" description:"端口号"`
	FetchedUnix int64  `json:"fetchedUnix,omitempty" description:"获取时间戳"`
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

type ClientListQueryReq struct {
	dto.PageQuery
	Hostname string `form:"hostname"`
	Status   int8   `form:"status"`
}
