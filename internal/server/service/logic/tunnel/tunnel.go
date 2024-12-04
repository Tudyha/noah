package tunnel

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
	"noah/pkg/mux"
	"strings"
	"sync"
)

type tunnel struct {
	id          uint
	tunnelType  uint8
	serverPort  int
	clientId    uint
	targetAddr  string
	closeSignal chan struct{}
	closeOnce   sync.Once
	service     *tunnelService
}

func newTunnel(id uint, tunnelType uint8, serverPort int, clientId uint, targetAddr string, s *tunnelService) *tunnel {
	return &tunnel{
		id:          id,
		tunnelType:  tunnelType,
		serverPort:  serverPort,
		clientId:    clientId,
		targetAddr:  targetAddr,
		closeSignal: make(chan struct{}),
		closeOnce:   sync.Once{},
		service:     s,
	}
}

func (t *tunnel) start() error {
	go t.listen()
	return nil
}

func (t *tunnel) stop() {
	t.closeOnce.Do(func() {
		close(t.closeSignal)
		// 向监听地址发送一个自连接以释放阻塞，fix 关闭端口监听后还可以接受一次连接
		if conn, err := net.Dial("tcp", fmt.Sprintf(":%d", t.serverPort)); err == nil {
			conn.Close()
		}
	})
}

func (t *tunnel) listen() {
	defer t.service.removeTunnel(t.id)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", t.serverPort))
	if err != nil {
		return
	}
	defer listener.Close()

	// 监听连接
	for {
		select {
		case <-t.closeSignal:
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				var opErr *net.OpError
				if errors.As(err, &opErr) && opErr.Err.Error() == "use of closed network connection" {
					return
				}
				continue
			}
			go t.handleConn(conn)
		}
	}
}

func (t *tunnel) handleConn(conn net.Conn) {
	defer conn.Close()

	network := "tcp"
	targetAddr := t.targetAddr

	var r *http.Request

	switch t.tunnelType {
	case 1:
	case 2:
		reader := bufio.NewReader(conn)
		req, err := http.ReadRequest(reader)
		if err != nil {
			return
		}
		targetAddr = req.Host
		if !strings.Contains(req.Host, ":") {
			if req.URL.Scheme == "https" {
				targetAddr = req.Host + ":443"
			}
			if req.URL.Scheme == "http" {
				targetAddr = req.Host + ":80"
			}
		}

		if req.Method == "CONNECT" {
			conn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
			req = nil
		}

		r = req
	}

	clientConn, err := t.service.newClientConn(t.clientId, network, targetAddr)
	if err != nil {
		return
	}
	defer clientConn.Close()
	if r != nil {
		r.Write(clientConn)
	}

	clientConn.(*mux.Conn).Copy(conn)
}
