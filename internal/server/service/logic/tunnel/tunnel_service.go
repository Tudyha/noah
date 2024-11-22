package tunnel

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"noah/internal/server/dao"
	"noah/internal/server/gateway"
	"noah/internal/server/middleware/log"
	"noah/internal/server/model"
	"noah/pkg/mux"
	"noah/pkg/request"
	"noah/pkg/response"
	"strconv"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/samber/do/v2"
)

type tunnelService struct {
	tunnelDao   dao.TunnelDao
	tunnelClose map[uint]chan struct{}
	gateway     *gateway.Gateway
}

func NewTunnelService(i do.Injector) (tunnelService, error) {
	s := tunnelService{
		tunnelDao:   do.MustInvoke[dao.TunnelDao](i),
		tunnelClose: make(map[uint]chan struct{}),
		gateway:     do.MustInvoke[*gateway.Gateway](i),
	}

	s.recoverTunnel()

	return s, nil
}

// NewTunnel 新建tunnel
func (c tunnelService) NewTunnel(id uint, tunnelReq request.CreateTunnelReq) (err error) {
	tunnelType := tunnelReq.TunnelType
	serverPort := tunnelReq.ServerPort
	clientIp := tunnelReq.ClientIp
	clientPort := tunnelReq.ClientPort

	tunnel := model.Tunnel{
		TunnelType: tunnelType,
		ClientId:   id,
		ClientIp:   clientIp,
		ClientPort: clientPort,
		ServerPort: serverPort,
	}
	tunnelId, err := c.tunnelDao.Save(tunnel)
	if err != nil {
		log.Error("Save tunnel error", map[string]interface{}{"clientId": id, "error": err})
		return err
	}

	// 服务端需要监听新端口
	go func() {
		err := c.listen(tunnelId)
		if err != nil {
			c.tunnelDao.UpdateStatus(tunnelId, 2, err.Error())
			return
		}
	}()

	return nil
}

func copy(src *mux.Conn, target io.ReadWriteCloser) {
	defer src.Close()
	defer target.Close()
	src.Copy(target)
}

// GetTunnelList 获取tunnel列表
func (c tunnelService) GetTunnelList(clientId uint) (res []response.GetTunnelListRes, err error) {
	list, err := c.tunnelDao.List(clientId)
	if err != nil {
		return nil, err
	}
	copier.Copy(&res, list)
	return res, nil
}

// DeleteTunnel 删除tunnel
func (c tunnelService) DeleteTunnel(id uint) (err error) {
	tunnel, err := c.tunnelDao.GetById(id)
	if err != nil {
		return err
	}
	err = c.tunnelDao.Delete(id)
	if err != nil {
		return err
	}

	// 关闭端口监听
	if _, ok := c.tunnelClose[id]; ok {
		close(c.tunnelClose[id])

		// 向监听地址发送一个自连接以释放阻塞，fix 关闭端口监听后还可以接受一次连接
		if conn, err := net.Dial("tcp", fmt.Sprintf(":%d", tunnel.ServerPort)); err == nil {
			conn.Close()
		}
	}
	return nil
}

// recoverTunnel 恢复tunnel
func (c tunnelService) recoverTunnel() {
	list, err := c.tunnelDao.List(0)
	if err != nil {
		return
	}

	for _, tunnel := range list {
		go func() {
			err := c.listen(tunnel.ID)
			if err != nil {
				c.tunnelDao.UpdateStatus(tunnel.ID, 2, err.Error())
				return
			}
		}()
	}
}

// NewTunnelConn 新建tunnel连接
func (c tunnelService) NewTunnelConn(clientId uint, tunnelType uint8, clientIp string, clientPort int) (*mux.Conn, error) {
	var network string
	switch tunnelType {
	case 1:
		network = "tcp"
	case 2:
		network = "tcp"
	}
	return c.gateway.NewClientConn(uint32(clientId), network, fmt.Sprintf("%s:%d", clientIp, clientPort))
}

func (c tunnelService) listen(tunnelId uint) error {
	tunnel, err := c.tunnelDao.GetById(tunnelId)
	if err != nil {
		log.Error("tunnel listen GetById error", map[string]interface{}{"tunnelId": tunnelId, "error": err})
		return err
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", tunnel.ServerPort))
	if err != nil {
		log.Error("listen error: "+err.Error(), nil)
		return err
	}

	c.tunnelClose[tunnelId] = make(chan struct{})

	c.tunnelDao.UpdateStatus(tunnelId, 1, "")

	// 监听连接
	for {
		select {
		case <-c.tunnelClose[tunnelId]:
			log.Info("tunnel listen close", map[string]interface{}{"tunnelId": tunnelId})
			delete(c.tunnelClose, tunnelId)
			listener.Close()
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "use of closed network connection" {
					// 监听器被关闭
					fmt.Println("Listener closed")
					return nil
				}
				continue
			}

			var rb []byte
			rb = nil
			if tunnel.TunnelType == 2 {
				_, addr, b, err, r := GetHost(conn)
				if err != nil {
					continue
				}
				if r.Method == "CONNECT" {
					conn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
					b = nil
				}
				tunnel.ClientIp = strings.Split(addr, ":")[0]
				tunnel.ClientPort, _ = strconv.Atoi(strings.Split(addr, ":")[1])
				rb = b
			}

			clientConn, err := c.NewTunnelConn(tunnel.ClientId, tunnel.TunnelType, tunnel.ClientIp, tunnel.ClientPort)
			if err != nil {
				continue
			}
			if rb != nil {
				_, err = clientConn.Write(rb)
				if err != nil {
					continue
				}
			}
			go copy(clientConn, conn)
		}
	}
}

func GetHost(conn net.Conn) (method, address string, rb []byte, err error, r *http.Request) {
	var b [32 * 1024]byte
	var n int
	if n, err = readRequest(conn, b[:]); err != nil {
		return
	}
	rb = b[:n]
	r, err = http.ReadRequest(bufio.NewReader(bytes.NewReader(rb)))
	if err != nil {
		return
	}
	hostPortURL, err := url.Parse(r.Host)
	if err != nil {
		address = r.Host
		err = nil
		return
	}
	if hostPortURL.Opaque == "443" {
		if strings.Index(r.Host, ":") == -1 {
			address = r.Host + ":443"
		} else {
			address = r.Host
		}
	} else {
		if strings.Index(r.Host, ":") == -1 {
			address = r.Host + ":80"
		} else {
			address = r.Host
		}
	}
	return
}

func readRequest(conn net.Conn, buf []byte) (n int, err error) {
	var rd int
	for {
		rd, err = conn.Read(buf[n:])
		if err != nil {
			return
		}
		n += rd
		if n < 4 {
			continue
		}
		if string(buf[n-4:n]) == "\r\n\r\n" {
			return
		}
		// buf is full, can't contain the request
		if n == cap(buf) {
			err = io.ErrUnexpectedEOF
			return
		}
	}
}
