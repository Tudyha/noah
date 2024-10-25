package conn

type Message struct {
	Flag   flag   `json:"flag"`
	ConnId uint32 `json:"connId"`
	Data   []byte `json:"data"`
	Msg    string `json:"msg"`
}

type flag int8

const (
	newConn flag = iota
	newConnOk
	data
	connClose
	ping
	pong
)

type LinkInfo struct {
	Network Network `json:"network"`
	Addr    string  `json:"addr"`
	CmdInfo CmdInfo `json:"data"`
}

type Network string

const (
	NetworkTcp Network = "tcp"
	NetworkCmd Network = "cmd"
	NetworkPty Network = "pty"
)

type CmdInfo struct {
	Cmd  string `json:"cmd"`
	Data []byte `json:"data"`
}

var (
	Command      = "command"
	Exit         = "exit"
	Update       = "update"
	Download     = "download"
	FileExplorer = "fileExplorer"
	SystemInfo   = "systemInfo"
)
