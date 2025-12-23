package io

import (
	"io"

	"github.com/gorilla/websocket"
)

type WebSocketReadWriteCloser struct {
	MessageType int
	Conn        *websocket.Conn
	textHandler func(reader io.Reader)
}

func (w *WebSocketReadWriteCloser) Read(p []byte) (n int, err error) {
	messageType, reader, err := w.Conn.NextReader()
	if err != nil {
		return
	}
	if messageType != w.MessageType {
		if w.textHandler != nil && messageType == websocket.TextMessage {
			w.textHandler(reader)
		}
		return 0, nil
	}
	return reader.Read(p)
}

func (w *WebSocketReadWriteCloser) SetTextHandler(h func(reader io.Reader)) {
	w.textHandler = h
}

func (w *WebSocketReadWriteCloser) Write(p []byte) (n int, err error) {
	return len(p), w.Conn.WriteMessage(w.MessageType, p)
}

func (w *WebSocketReadWriteCloser) Close() error {
	return w.Conn.Close()
}
