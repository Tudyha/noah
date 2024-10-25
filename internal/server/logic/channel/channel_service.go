package channel

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"noah/internal/server/dao"
	"noah/internal/server/enum"
	"noah/internal/server/gateway"
	"noah/internal/server/middleware/log"
	"noah/internal/server/request"
	"noah/internal/server/response"
	"noah/pkg/conn"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/jinzhu/copier"
	"github.com/samber/do/v2"
)

type Service struct {
	channelClose map[uint]chan struct{} // channel关闭
	gateway      *gateway.Gateway
}

// NewChannelService 初始化
func NewChannelService(i do.Injector) (*Service, error) {
	s := &Service{
		channelClose: make(map[uint]chan struct{}),
		gateway:      do.MustInvoke[*gateway.Gateway](i),
	}

	// 恢复channel fixme 用到service才初始化，channel恢复不及时
	s.recoverChannel()

	return s, nil
}

// NewChannel 新建channel
func (c Service) NewChannel(id uint, channelReq request.CreateChannelReq, wsconn *websocket.Conn) (err error) {
	channelType := channelReq.ChannelType
	serverPort := channelReq.ServerPort
	clientIp := channelReq.ClientIp
	clientPort := channelReq.ClientPort

	if channelType == enum.Pty {
		// pty模式，web -> websocket -> server -> websocket -> client
		clientConn, err := c.NewChannelConn(id, enum.Pty, clientIp, clientPort)
		if err != nil {
			return err
		}
		target := &conn.WebSocketReaderWriterCloser{Conn: wsconn}
		go copy(clientConn, target)
		return nil
	}

	// channel配置信息写进数据库，服务重启后可以从数据库恢复
	channel := dao.Channel{
		ChannelType: channelType,
		ClientId:    id,
		ClientIp:    clientIp,
		ClientPort:  clientPort,
		ServerPort:  serverPort,
	}
	channelId, err := dao.GetChannelDao().Save(channel)
	if err != nil {
		log.Error("Save channel error", map[string]interface{}{"clientId": id, "error": err})
		return err
	}

	if channelType == enum.Tcp || channelType == enum.Http {
		// 服务端需要监听新端口
		go func() {
			err := c.listen(channelId)
			if err != nil {
				dao.GetChannelDao().UpdateStatus(channelId, enum.ChannelStatusDisconnected, err.Error())
				return
			}
		}()
	}

	return nil
}

func copy(src *conn.Conn, target io.ReadWriteCloser) {
	defer src.Close()
	defer target.Close()
	src.Copy(target)
}

// GetChannelList 获取channel列表
func (c Service) GetChannelList(clientId uint) (res []response.GetChannelListRes, err error) {
	list, err := dao.GetChannelDao().List(clientId)
	if err != nil {
		return nil, err
	}
	copier.Copy(&res, list)
	return res, nil
}

// DeleteChannel 删除channel
func (c Service) DeleteChannel(id uint) (err error) {
	channel, err := dao.GetChannelDao().GetById(id)
	if err != nil {
		return err
	}
	err = dao.GetChannelDao().Delete(id)
	if err != nil {
		return err
	}

	// 关闭端口监听
	if _, ok := c.channelClose[id]; ok {
		close(c.channelClose[id])

		// 向监听地址发送一个自连接以释放阻塞，fix 关闭端口监听后还可以接受一次连接
		if conn, err := net.Dial("tcp", fmt.Sprintf(":%d", channel.ServerPort)); err == nil {
			conn.Close()
		}
	}
	return nil
}

// recoverChannel 恢复channel
func (c Service) recoverChannel() {
	list, err := dao.GetChannelDao().List(0)
	if err != nil {
		return
	}

	for _, channel := range list {
		if (channel.ChannelType == enum.Tcp || channel.ChannelType == enum.Http) && channel.Status == enum.ChannelStatusConnected {
			go func() {
				err := c.listen(channel.ID)
				if err != nil {
					dao.GetChannelDao().UpdateStatus(channel.ID, enum.ChannelStatusDisconnected, err.Error())
					return
				}
			}()
		}
	}
}

// NewChannelConn 新建channel连接
func (c Service) NewChannelConn(clientId uint, channelType enum.ChannelType, clientIp string, clientPort int) (*conn.Conn, error) {
	var network conn.Network
	switch channelType {
	case enum.Pty:
		network = conn.NetworkPty
	case enum.Tcp:
		network = conn.NetworkTcp
	case enum.Http:
		network = conn.NetworkTcp

	}
	return c.gateway.NewClientConn(clientId, network, fmt.Sprintf("%s:%d", clientIp, clientPort), conn.CmdInfo{})
}

func (c Service) listen(channelId uint) error {
	channel, err := dao.GetChannelDao().GetById(channelId)
	if err != nil {
		log.Error("channel listen GetById error", map[string]interface{}{"channelId": channelId, "error": err})
		return err
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", channel.ServerPort))
	if err != nil {
		log.Error("listen error: "+err.Error(), nil)
		return err
	}

	c.channelClose[channelId] = make(chan struct{})

	dao.GetChannelDao().UpdateStatus(channelId, enum.ChannelStatusConnected, "")

	// 监听连接
	for {
		select {
		case <-c.channelClose[channelId]:
			log.Info("channel listen close", map[string]interface{}{"channelId": channelId})
			delete(c.channelClose, channelId)
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
			if channel.ChannelType == enum.Http {
				_, addr, b, err, r := GetHost(conn)
				if err != nil {
					continue
				}
				if r.Method == "CONNECT" {
					conn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
					b = nil
				}
				if channel.ChannelType == enum.Http {
					channel.ClientIp = strings.Split(addr, ":")[0]
					channel.ClientPort, _ = strconv.Atoi(strings.Split(addr, ":")[1])
				}
				rb = b
			}

			clientConn, err := c.NewChannelConn(channel.ClientId, enum.Tcp, channel.ClientIp, channel.ClientPort)
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
