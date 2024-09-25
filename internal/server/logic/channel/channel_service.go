package channel

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/copier"
	"net"
	"noah/internal/server/dao"
	"noah/internal/server/enum"
	"noah/internal/server/gateway"
	"noah/internal/server/middleware/log"
	"noah/internal/server/request"
	"noah/internal/server/response"
	"noah/internal/server/utils"
	"sync"
)

type Service struct {
	mu           *sync.Mutex
	channelConns map[string]*Conn       // channel连接
	channelClose map[uint]chan struct{} // channel关闭
	gateway      *gateway.Gateway
}

func NewChannelService(gateway *gateway.Gateway) *Service {
	s := &Service{
		mu:           &sync.Mutex{},
		channelConns: make(map[string]*Conn),
		channelClose: make(map[uint]chan struct{}),
		gateway:      gateway,
	}

	// 恢复channel
	s.recoverChannel()

	// health check
	//go func() {
	//	for {
	//		log.Info("channel health check", map[string]interface{}{
	//			"channelConns": s.channelConns,
	//		})
	//
	//		time.Sleep(time.Second * 10)
	//	}
	//}()

	return s
}

type Conn struct {
	clientId uint       // 客户端id
	connId   string     // 连接id
	conn     Connection // 连接
}

// NewChannel 新建channel
func (c Service) NewChannel(id uint, channelReq request.CreateChannelReq, conn *websocket.Conn) (err error) {
	channelType := channelReq.ChannelType
	serverPort := channelReq.ServerPort
	clientIp := channelReq.ClientIp
	clientPort := channelReq.ClientPort

	if channelType == enum.Pty {
		channel, err := c.NewChannelConn(id, &WebSocketConnection{conn: conn}, enum.Pty, clientIp, clientPort)
		if err != nil {
			log.Error("NewChannelConn error", map[string]interface{}{"clientId": id, "error": err})
			return err
		}
		go channel.read(c)
		go channel.write(c)
		return nil
	}

	// channel配置信息写进数据库，服务重启后可以从数据恢复
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

	if channelType == enum.Tcp {
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
		if channel.ChannelType == enum.Tcp && channel.Status == enum.ChannelStatusConnected {
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
func (c Service) NewChannelConn(id uint, conn Connection, channelType enum.ChannelType, clientIp string, clientPort int) (*Conn, error) {
	channelConn := &Conn{
		clientId: id,
		connId:   utils.RandString(16),
		conn:     conn,
	}

	// 通知客户端打开对应通道
	_, err := c.gateway.ClientWebsocketWrite(id, channelConn.connId, enum.MessageTypeChannel, &request.ChannelRequest{
		Action:      "open",
		ChannelId:   channelConn.connId,
		ChannelType: channelType,
		ChannelData: nil,
		LocalIp:     clientIp,
		LocalPort:   clientPort,
	})
	if err != nil {
		log.Error("NewChannelConn open error", map[string]interface{}{"clientId": id, "error": err})
		return nil, err
	}

	c.channelConns[channelConn.connId] = channelConn
	return channelConn, nil
}

func (c Service) closeChannelConn(connId string) error {
	conn, ok := c.channelConns[connId]
	if !ok {
		return errors.New("channel conn not found")
	}

	if conn.conn != nil {
		err := conn.conn.Close()
		if err != nil {
			return err
		}
	}

	delete(c.channelConns, conn.connId)

	return nil
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

			channelCoon, err := c.NewChannelConn(channel.ClientId, &TCPConnection{conn: conn}, enum.Tcp, channel.ClientIp, channel.ClientPort)
			if err != nil {
				log.Error("NewChannelConn error", map[string]interface{}{"error": err})
				continue
			}

			go channelCoon.read(c)
			go channelCoon.write(c)
		}
	}
}

// read 从连接中读取数据，并转发给客户端
func (conn *Conn) read(c Service) {
	defer func() {
		err := c.closeChannelConn(conn.connId)
		if err != nil {
			log.Error("closeChannelConn error", map[string]interface{}{"error": err})
			return
		}
	}()
	for {
		buffer, err := conn.conn.ReadMessage()
		if err != nil {
			return
		}

		wsData := request.ChannelRequest{
			Action:      "write",
			ChannelId:   conn.connId,
			ChannelData: buffer,
		}
		_, err = c.gateway.ClientWebsocketWrite(conn.clientId, conn.connId, enum.MessageTypeChannel, wsData)
		if err != nil {
			return
		}
	}
}

// write 从客户端读取数据，并转发给连接
func (conn *Conn) write(c Service) {
	defer func() {
		c.gateway.UnSubscribeMessage(conn.clientId, conn.connId)

		err := c.closeChannelConn(conn.connId)
		if err != nil {
			return
		}
	}()

	ch := make(chan request.Message, 32)

	c.gateway.SubscribeMessage(conn.clientId, conn.connId, func(message request.Message) {
		ch <- message
	})

	for message := range ch {
		var channelRequest request.ChannelRequest
		err := json.Unmarshal(message.Data, &channelRequest)
		if err != nil {
			continue
		}

		if channelRequest.ChannelId != conn.connId {
			log.Warn("channelId not match", nil)
			continue
		}

		if channelRequest.Action != "write" {
			log.Warn("action not match", nil)
			continue
		}

		if err := conn.conn.WriteMessage(channelRequest.ChannelData); err != nil {
			log.Error("Error writing to TCP:"+err.Error(), nil)
			break
		}
	}
	close(ch)
}
