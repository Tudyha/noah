package enum

type ChannelType int

const (
	UnknownChannelType ChannelType = iota
	Pty
	Tcp
)

type ChannelStatus int8

const (
	ChannelStatusWait ChannelStatus = iota
	ChannelStatusConnected
	ChannelStatusDisconnected
)
