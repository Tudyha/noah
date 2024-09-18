package channel

import (
	"github.com/gorilla/websocket"
	"net"
)

// Connection 接口统一了 net.Conn 和 websocket.Conn 的行为
type Connection interface {
	ReadMessage() ([]byte, error)
	WriteMessage([]byte) error
	Close() error
}

// TCPConnection 包装了 net.Conn 来实现 Connection 接口
type TCPConnection struct {
	conn net.Conn
}

func (t *TCPConnection) ReadMessage() ([]byte, error) {
	buf := make([]byte, 1024)
	n, err := t.conn.Read(buf)
	return buf[:n], err
}

func (t *TCPConnection) WriteMessage(message []byte) error {
	_, err := t.conn.Write(message)
	return err
}

func (t *TCPConnection) Close() error {
	return t.conn.Close()
}

// WebSocketConnection 包装了 websocket.Conn 来实现 Connection 接口
type WebSocketConnection struct {
	conn *websocket.Conn
}

func (w *WebSocketConnection) ReadMessage() ([]byte, error) {
	_, message, err := w.conn.ReadMessage()
	return message, err
}

func (w *WebSocketConnection) WriteMessage(message []byte) error {
	return w.conn.WriteMessage(websocket.TextMessage, message)
}

func (w *WebSocketConnection) Close() error {
	return w.conn.Close()
}
