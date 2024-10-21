package entitie

type MessageType int

const (
	MessageTypeUnknown MessageType = iota
	MessageTypeCommand
	MessageTypePty
	MessageTypeDownload
	MessageTypeUpdate
	MessageTypeExit
	MessageTypeFileExplorer
	MessageTypeSystemInfo
	MessageTypeChannel
)

func MessageTypeFromString(s string) MessageType {
	switch s {
	case "1":
		return MessageTypeCommand
	case "3":
		return MessageTypeDownload
	case "4":
		return MessageTypeUpdate
	case "5":
		return MessageTypeExit
	case "6":
		return MessageTypeFileExplorer
	case "7":
		return MessageTypeSystemInfo
	}
	return MessageTypeUnknown
}

type Message struct {
	MessageId   string      `json:"messageId"`
	MessageType MessageType `json:"messageType,omitempty"`
	Data        []byte      `json:"data,omitempty"`
	Error       string      `json:"error,omitempty"`
}

type CommandReq struct {
	Command string `json:"command,omitempty"`
}

type ChannelReq struct {
	Action      string `json:"action,omitempty"`
	ChannelId   string `json:"channelId,omitempty"`
	ChannelType int    `json:"channelType,omitempty"`
	ChannelData []byte `json:"channelData,omitempty"`
	LocalIp     string `json:"localIp,omitempty"`
	LocalPort   int    `json:"localPort,omitempty"`
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
