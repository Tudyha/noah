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
	Network string `json:"network"`
	Addr    string `json:"addr"`
}

type ptyData struct {
	Type  string `json:"type"`
	Data  any    `json:"data"`
	High  int    `json:"high"`
	Width int    `json:"width"`
}
