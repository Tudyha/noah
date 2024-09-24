package enum

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
