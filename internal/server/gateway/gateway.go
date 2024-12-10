package gateway

import (
	"fmt"
	"net"
	"noah/pkg/mux"
	"noah/pkg/mux/message"
	"sync"

	"noah/pkg/utils"

	"github.com/samber/do/v2"
)

type Gateway struct {
	conns       sync.Map
	pongHandler func(uint32, []byte)
}

func NewGateway(i do.Injector) (Gateway, error) {
	return Gateway{
		conns: sync.Map{},
	}, nil
}

func (g *Gateway) SetPongHandler(h func(uint32, []byte)) {
	g.pongHandler = h
}

func (g *Gateway) HanderConn(clientId uint32, conn net.Conn) {
	if old, ok := g.conns.Load(clientId); ok {
		old.(*mux.Mux).Close()
	}
	m := mux.NewMux(conn, conn)
	err := m.Start()
	if err != nil {
		return
	}

	g.conns.Store(clientId, m)
	m.SetClosedCallbackHandler(func() {
		g.conns.Delete(clientId)
		conn.Close()
	})
	m.SetPongHandler(func(data []byte) {
		if g.pongHandler != nil {
			g.pongHandler(clientId, data)
		}
	})
}

func (g *Gateway) NewClientConn(clientId uint32, network, addr string) (*mux.Conn, error) {
	if v, ok := g.conns.Load(clientId); !ok {
		return nil, fmt.Errorf("client not found")
	} else {
		if m, ok := v.(*mux.Mux); !ok {
			return nil, fmt.Errorf("client not found")
		} else {
			return m.Dial(network, addr)
		}
	}
}

func (g *Gateway) SendCommand(clientId uint, cmd string, data any, needResult bool) ([]byte, error) {
	cd, err := utils.AnyToJsonBytes(data)
	if err != nil {
		return nil, err
	}
	cmdInfo := message.CmdInfo{
		Cmd:  cmd,
		Data: cd,
	}
	addr, err := utils.AnyToJsonBytes(cmdInfo)
	if err != nil {
		return nil, err
	}

	conn, err := g.NewClientConn(uint32(clientId), "cmd", string(addr))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if !needResult {
		return nil, err
	}
	res, err := conn.ReadFull()
	if err != nil {
		return nil, err
	}
	return res, err
}
