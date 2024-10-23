package gateway

import (
	"errors"
	"fmt"
	"noah/internal/server/enum"
	"noah/pkg/conn"
	"noah/pkg/utils"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/samber/do/v2"
)

type Gateway struct {
	mu      *sync.Mutex
	clients map[uint]*conn.Mux
}

func NewGateway(i do.Injector) (*Gateway, error) {
	g := &Gateway{
		mu:      &sync.Mutex{},
		clients: make(map[uint]*conn.Mux),
	}
	return g, nil
}

// NewClientWebsocketConn 新增ws连接
func (g Gateway) NewClientWebsocketConn(clientId uint, connection *websocket.Conn) error {
	// 验证连接是否有效
	if connection == nil || connection.RemoteAddr() == nil {
		return errors.New("invalid connection")
	}
	g.mu.Lock()
	defer g.mu.Unlock()

	// 更新或添加新连接
	m := conn.NewMux(clientId, connection)
	g.clients[clientId] = m

	return nil
}

func (g Gateway) NewClientConn(clientId uint, network string, addr string) (*conn.Conn, error) {
	if client, ok := g.clients[clientId]; ok {
		if client.Closed {
			delete(g.clients, clientId)
			return nil, errors.New("client connection closed")
		}
		return client.NewConn(network, addr)
	} else {
		return nil, errors.New("client not found")
	}
}

// SendCommand 执行命令
func (g Gateway) SendCommand(clientId uint, messageType enum.MessageType, data any, needResult bool) (string, error) {
	addr := rune(messageType)
	srcConn, err := g.NewClientConn(clientId, "cmd", fmt.Sprintf("%d", addr))
	if err != nil {
		return "", err
	}
	defer srcConn.Close()

	b, err := utils.AnyToBytes(data)
	if err != nil {
		return "", err
	}
	_, err = srcConn.Write(b)
	if err != nil {
		return "", err
	}
	if !needResult {
		return "", nil
	}
	res, err := srcConn.ReadFull()
	if err != nil {
		return "", err
	}
	return string(res), err
}
