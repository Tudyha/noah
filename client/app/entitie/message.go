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
	MessageTypeProcess
	MessageTypeChannel
)

type Message struct {
	MessageId   uint64      `json:"messageId"`
	MessageType MessageType `json:"messageType,omitempty"`
	Data        []byte      `json:"data,omitempty"`
	Error       string      `json:"error,omitempty"`
}

type CommandRequest struct {
	Command string `json:"command,omitempty"`
}

type PtyRequest struct {
	Action      string `json:"action,omitempty"`
	ChannelId   string `json:"channelId,omitempty"`
	ChannelData []byte `json:"channelData,omitempty"`
}

type DownloadRequest struct {
	Filename string `json:"filename"`
	Path     string `json:"path"`
}
