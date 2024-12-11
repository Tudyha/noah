package tunnel

import (
	"errors"
	"fmt"
	"github.com/shadowsocks/go-shadowsocks2/core"
	"github.com/shadowsocks/go-shadowsocks2/socks"
	"io"
	"log"
	"net"
	"os"
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
	var key []byte
	cipher := "AES-256-GCM"
	password := "123456"

	ciph, err := core.PickCipher(cipher, key, password)
	if err != nil {
		log.Fatal(err)
	}

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
			go t.handleConn(conn, ciph.StreamConn)
		}
	}
}

func (t *tunnel) handleConn(c net.Conn, shadow func(net.Conn) net.Conn) {
	defer c.Close()

	network := "tcp"
	targetAddr := t.targetAddr

	var sc net.Conn

	switch t.tunnelType {
	case 1:
		sc = c
	case 2:
		sc = shadow(c)

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

	log.Printf("proxy %s <-> %s", c.RemoteAddr(), targetAddr)
	if err = relay(sc, rc); err != nil {
		log.Printf("relay error: %v", err)
	}
}

// relay copies between left and right bidirectionally
func relay(left, right net.Conn) error {
	var err, err1 error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err1 = io.Copy(right, left)
	}()
	_, err = io.Copy(left, right)
	wg.Wait()
	if err1 != nil && !errors.Is(err1, os.ErrDeadlineExceeded) { // requires Go 1.15+
		return err1
	}
	if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
		return err
	}
	return nil
}
