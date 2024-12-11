package tunnel

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	myio "noah/pkg/io"
	"sync"

	"github.com/shadowsocks/go-shadowsocks2/core"
	"github.com/shadowsocks/go-shadowsocks2/socks"
)

type tunnel struct {
	id          uint          // tunnel id
	tunnelType  uint8         // 隧道类型: 1.tcp 2.shadowsocks
	serverPort  int           // 服务端监听端口
	clientId    uint          // 客户端id
	targetAddr  string        // 目标地址
	closeSignal chan struct{} // 关闭信号
	closeOnce   sync.Once
	service     *tunnelService          // tunnel service
	shadow      func(net.Conn) net.Conn // shadowsocks 加密
}

// cipher: shadowsocks 加密方式: CHACHA20-IETF-POLY1305、AES-128-GCM、AES-256-GCM
// password: shadowsocks 密码
func newTunnel(id uint, tunnelType uint8, serverPort int, clientId uint, targetAddr string, cipher, password string, s *tunnelService) (*tunnel, error) {
	var shadow func(net.Conn) net.Conn
	if tunnelType == 2 {
		var key []byte

		ciph, err := core.PickCipher(cipher, key, password)
		if err != nil {
			return nil, err
		}
		shadow = ciph.StreamConn
	}
	return &tunnel{
		id:          id,
		tunnelType:  tunnelType,
		serverPort:  serverPort,
		clientId:    clientId,
		targetAddr:  targetAddr,
		closeSignal: make(chan struct{}),
		closeOnce:   sync.Once{},
		service:     s,
		shadow:      shadow,
	}, nil
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

func (t *tunnel) handleConn(c net.Conn) {
	defer c.Close()

	network := "tcp"
	targetAddr := t.targetAddr

	var sc net.Conn

	switch t.tunnelType {
	case 1:
		sc = c
	case 2:
		sc = t.shadow(c)

		tgt, err := socks.ReadAddr(sc)
		if err != nil {
			log.Printf("failed to get target address from %v: %v", c.RemoteAddr(), err)
			// drain c to avoid leaking server behavioral features
			// see https://www.ndss-symposium.org/ndss-paper/detecting-probe-resistant-proxies/
			_, err = io.Copy(io.Discard, c)
			if err != nil {
				log.Printf("discard error: %v", err)
			}
			return
		}

		targetAddr = tgt.String()
	}

	rc, err := t.service.newClientConn(t.clientId, network, targetAddr)
	if err != nil {
		log.Printf("failed to connect to %v: %v", targetAddr, err)
		return
	}
	defer rc.Close()

	if err = myio.Copy(sc, rc); err != nil {
		log.Printf("relay error: %v", err)
	}
}
