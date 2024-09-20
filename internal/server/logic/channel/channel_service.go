package channel

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"net"
	"noah/internal/server/dao"
	"noah/internal/server/enum"
	"noah/internal/server/middleware/log"
	"noah/internal/server/request"
	"noah/internal/server/response"
	"noah/internal/server/utils"
	"sync"
	"time"

	"github.com/duke-git/lancet/v2/slice"

	"github.com/gorilla/websocket"
)

var (
	messageId = 0
)

type Service struct {
	mu            *sync.Mutex
	clients       map[uint]*websocket.Conn        //客户端命令执行websocket
	messageResult map[uint64]chan request.Message //命令执行结果
	channelConns  map[string]*Conn                // channel连接
	messageMq     map[string]chan []byte          //channel消息队列
	channelClose  map[uint]chan struct{}          // channel关闭
}

func NewChannelService() *Service {
	s := &Service{
		mu:            &sync.Mutex{},
		clients:       make(map[uint]*websocket.Conn),
		messageResult: make(map[uint64]chan request.Message),
		channelConns:  make(map[string]*Conn),
		messageMq:     make(map[string]chan []byte),
		channelClose:  make(map[uint]chan struct{}),
	}

	// 恢复channel
	s.recoverChannel()

	return s
}

// NewClientWebsocketConn 新增ws连接
func (c Service) NewClientWebsocketConn(id uint, connection *websocket.Conn) error {
	// 验证连接是否有效
	if connection == nil || connection.RemoteAddr() == nil {
		return errors.New("invalid connection")
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	// 更新或添加新连接
	c.clients[id] = connection

	// 启动一个goroutine用于读取客户端的websocket消息
	go c.clientWebsocketRead(id)
	return nil
}

// getClientWebsocketConn 获取客户端websocket连接
func (c Service) getClientWebsocketConn(clientID uint) (*websocket.Conn, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	conn, found := c.clients[clientID]
	if !found {
		return nil, errors.New("client not found")
	}
	return conn, nil
}

// Exit 关闭连接
func (c Service) Exit(id uint) error {
	if err := c.removeClientWebsocketConnection(id); err != nil {
		return err
	}
	return nil
}

// removeClientWebsocketConnection 删除客户端websocket连接
func (c Service) removeClientWebsocketConnection(clientID uint) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if conn, found := c.clients[clientID]; !found {
	} else {
		err := conn.Close()
		if err != nil {
		}
	}
	delete(c.clients, clientID)
	return nil
}

// clientWebsocketWrite 向发送消息
func (c Service) clientWebsocketWrite(id uint, websocketMessageType int, messageType enum.MessageType, data any, errMsg string) (messageId uint64, err error) {
	client, err := c.getClientWebsocketConn(id)
	if err != nil {
		return 0, err
	}

	var messageByteData []byte
	switch data := data.(type) {
	case []byte:
		messageByteData = data
	case string:
		messageByteData = []byte(data)
	default:
		messageByteData, err = json.Marshal(data)
		if err != nil {
			return 0, err
		}
	}

	messageId = messageId + 1

	req := request.Message{
		MessageId:   messageId,
		MessageType: messageType,
		Data:        messageByteData,
		Error:       errMsg,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return 0, err
	}

	err = client.WriteMessage(websocketMessageType, body)
	if err != nil {
		log.Error("write client websocket message error, websocket exit", map[string]interface{}{"clientId": id, "error": err})
		c.Exit(id)
		return 0, err
	}

	return messageId, nil
}

// clientWebsocketRead 读取客户端的websocket消息
func (c Service) clientWebsocketRead(id uint) {
	defer func() {
		log.Info("client websocket read fail, websocket exit", map[string]interface{}{"clientId": id})
		c.Exit(id)
	}()
	client, err := c.getClientWebsocketConn(id)
	if err != nil {
		return
	}
	for {
		_, wsMessage, err := client.ReadMessage()
		if err != nil {
			log.Error("read message error", map[string]interface{}{"clientId": id, "error": err})
			return
		}
		var message request.Message
		if err := json.Unmarshal(wsMessage, &message); err != nil {
			log.Error("parse message error", map[string]interface{}{"clientId": id, "error": err})
			continue
		}

		c.handleMessage(id, message)
	}
}

// handleMessage 处理客户端发送的消息
func (c Service) handleMessage(id uint, message request.Message) {
	defer func() {
		if isNeedResult(message.MessageType) {
			//判断chan是否已初始化以及是否已关闭
			if _, ok := c.messageResult[message.MessageId]; ok {
				c.messageResult[message.MessageId] <- message
			}
		}
	}()

	if message.Error != "" {
		log.Warn("receive message error", map[string]interface{}{"clientId": id, "error": message.Error})
		return
	}

	switch message.MessageType {
	case enum.MessageTypeChannel:
		var channelRequest request.ChannelRequest
		err := json.Unmarshal(message.Data, &channelRequest)
		if err != nil {
			log.Error("parse pty request error", map[string]interface{}{"clientId": id, "error": err})
			return
		}
		switch channelRequest.Action {
		case "write":
			_, ok := c.channelConns[channelRequest.ChannelId]
			if !ok {
				break
			}

			// 将数据写入channel
			c.messageMq[channelRequest.ChannelId] <- channelRequest.ChannelData
		}
	default:

	}
}

func isNeedResult(messageType enum.MessageType) bool {
	return slice.Contain([]enum.MessageType{
		enum.MessageTypeCommand,
		enum.MessageTypeProcess,
		enum.MessageTypeFileExplorer}, messageType)
}

// SendCommand 执行命令
func (c Service) SendCommand(id uint, messageType enum.MessageType, data any) (string, error) {
	msgId, err := c.clientWebsocketWrite(id, websocket.TextMessage, messageType, data, "")
	if err != nil {
		return "", err
	}

	if !isNeedResult(messageType) {
		// 不需要命令执行结果，直接返回
		return "", nil
	}

	c.messageResult[msgId] = make(chan request.Message, 1)
	defer func() {
		if _, ok := c.messageResult[msgId]; ok {
			close(c.messageResult[msgId])
			delete(c.messageResult, msgId)
		}
	}()

	// 等待结果
	var result request.Message
	select {
	case result = <-c.messageResult[msgId]:

	case <-time.After(5 * time.Second):
		return "", errors.New("command exec timeout")
	}

	if result.Error != "" {
		return "", errors.New(result.Error)
	}

	return utils.ByteToString(result.Data), nil
}

type Conn struct {
	clientId uint       // 客户端id
	connId   string     // 连接id
	conn     Connection // 连接
}

func (c Service) NewChannel(id uint, channelReq request.CreateChannelReq, conn *websocket.Conn) (err error) {
	_, err = c.getClientWebsocketConn(id)
	if err != nil {
		log.Error("client websocket connection not found", map[string]interface{}{"clientId": id, "error": err})
		return err
	}

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

func (c Service) GetChannelList(clientId uint) (res []response.GetChannelListRes, err error) {
	list, err := dao.GetChannelDao().List(clientId)
	if err != nil {
		return nil, err
	}
	copier.Copy(&res, list)
	return res, nil
}

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

func (c Service) NewChannelConn(id uint, conn Connection, channelType enum.ChannelType, clientIp string, clientPort int) (*Conn, error) {
	channelConn := &Conn{
		clientId: id,
		connId:   utils.RandString(16),
		conn:     conn,
	}

	// 通知客户端打开对应通道
	_, err := c.SendCommand(id, enum.MessageTypeChannel, &request.ChannelRequest{
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
	c.messageMq[channelConn.connId] = make(chan []byte, 32)
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
	if _, ok := c.messageMq[conn.connId]; ok {
		close(c.messageMq[conn.connId])
		delete(c.messageMq, conn.connId)
	}

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
		_, err = c.clientWebsocketWrite(conn.clientId, websocket.TextMessage, enum.MessageTypeChannel, wsData, "")
		if err != nil {
			return
		}
	}
}

func (conn *Conn) write(c Service) {
	defer func() {
		err := c.closeChannelConn(conn.connId)
		if err != nil {
			log.Error("closeChannelConn error", map[string]interface{}{"error": err})
			return
		}
	}()
	for {
		data := <-c.messageMq[conn.connId]
		if err := conn.conn.WriteMessage(data); err != nil {
			log.Error("Error writing to TCP:"+err.Error(), nil)
			return
		}
	}
}
