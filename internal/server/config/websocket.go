package config

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var (
	maxMessageSize = 512
	Upgrader       = websocket.Upgrader{
		ReadBufferSize:  maxMessageSize,
		WriteBufferSize: maxMessageSize,
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// MessageWait Time allowed to write or read a message.
	MessageWait = 10 * time.Second
)
