package gateway

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/samber/do/v2"
	"net/http"
	"noah/pkg/conn"
	"noah/pkg/utils"
	"sync"
)

var (
	maxMessageSize = 32 * 1024
	upgrader       = websocket.Upgrader{
		ReadBufferSize:  maxMessageSize,
		WriteBufferSize: maxMessageSize,
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Gateway struct {
	clients sync.Map
}

func NewGateway(i do.Injector) (*Gateway, error) {
	g := &Gateway{
		clients: sync.Map{},
	}
	return g, nil
}

func UpgradeWebsocket(ctx *gin.Context) (*websocket.Conn, error) {
	return upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
}

// NewClientWebsocketConn 新增ws连接
func (g *Gateway) NewClientWebsocketConn(clientId uint, ctx *gin.Context) error {
	connection, err := UpgradeWebsocket(ctx)
	if err != nil {
		return errors.New("upgrade error")
	}

	// 删除旧的连接
	g.DelClient(clientId)

	// 更新或添加新连接
	m := conn.NewMux(clientId, connection)
	g.clients.Store(clientId, m)

	return nil
}

func (g *Gateway) DelClient(clientId uint) {
	if client, ok := g.clients.Load(clientId); ok {
		mux := client.(*conn.Mux)
		mux.Close()

		g.clients.Delete(clientId)
	}
}

// NewClientConn 新建客户端连接
func (g *Gateway) NewClientConn(clientId uint, network conn.Network, addr string, cmd conn.CmdInfo) (*conn.Conn, error) {
	if client, ok := g.clients.Load(clientId); ok {
		client := client.(*conn.Mux)
		if client.Closed {
			g.DelClient(clientId)
			return nil, errors.New("client disconnect")
		}
		return client.NewConn(network, addr, cmd)
	} else {
		return nil, errors.New("client not found")
	}
}

// SendCommand 执行命令
func (g *Gateway) SendCommand(clientId uint, cmd string, data any, needResult bool) (string, error) {
	b, err := utils.AnyToBytes(data)
	if err != nil {
		return "", err
	}
	srcConn, err := g.NewClientConn(clientId, conn.NetworkCmd, "", conn.CmdInfo{Cmd: cmd, Data: b})
	if err != nil {
		return "", err
	}
	defer srcConn.Close()

	if !needResult {
		return "", nil
	}
	res, err := srcConn.ReadFull()
	if err != nil {
		return "", err
	}
	return string(res), err
}
