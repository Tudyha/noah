package request

import "noah/internal/server/enum"

type Message struct {
	MessageId   uint64           `json:"messageId"`
	MessageType enum.MessageType `json:"messageType,omitempty"`
	Data        []byte           `json:"data,omitempty"`
	Error       string           `json:"error,omitempty"`
}

type CommandRequest struct {
	Command string `json:"command,omitempty"`
}

type ChannelRequest struct {
	Action      string           `json:"action,omitempty"`
	ChannelId   string           `json:"channelId,omitempty"`
	ChannelType enum.ChannelType `json:"channelType,omitempty"`
	ChannelData []byte           `json:"channelData,omitempty"`
	Addr        string           `json:"addr,omitempty"`
}

type DownloadRequest struct {
	Filename string `json:"filename"`
	Path     string `json:"path"`
}
