package channel

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"noah/internal/server/enum"
	"noah/internal/server/middleware/log"
	"noah/internal/server/request"
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
	clients       map[uint]*websocket.Conn //客户端命令执行websocket
	messageResult map[uint64]chan request.Message
	channels      map[string]*Channel
	messageMq     map[string]chan []byte
}

func NewChannelService() *Service {
	return &Service{
		mu:            &sync.Mutex{},
		clients:       make(map[uint]*websocket.Conn),
		messageResult: make(map[uint64]chan request.Message),
		channels:      make(map[string]*Channel),
		messageMq:     make(map[string]chan []byte),
	}
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

// clientWebsocketWrite 发送消息
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
		if message.Error != "" {
			log.Warn("receive message error", map[string]interface{}{"clientId": id, "error": err})
			continue
		}
		c.handleMessage(id, message)
	}
}

// handleMessage 处理客户端发送的消息
func (c Service) handleMessage(id uint, message request.Message) {
	switch message.MessageType {
	case enum.MessageTypeChannel:
		var ptyRequest request.ChannelRequest
		err := json.Unmarshal(message.Data, &ptyRequest)
		if err != nil {
			log.Error("parse pty request error", map[string]interface{}{"clientId": id, "error": err})
			return
		}
		switch ptyRequest.Action {
		case "write":
			_, ok := c.channels[ptyRequest.ChannelId]
			if !ok {
				break
			}
			c.messageMq[ptyRequest.ChannelId] <- ptyRequest.ChannelData
		}
	default:

	}

	if isNeedResult(message.MessageType) {
		//判断chan是否已初始化以及是否已关闭
		if _, ok := c.messageResult[message.MessageId]; ok {
			c.messageResult[message.MessageId] <- message
		}
	}
}

func isNeedResult(messageType enum.MessageType) bool {
	if messageType == enum.MessageTypeChannel {

	}
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

	c.messageResult[msgId] = make(chan request.Message)
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

type Channel struct {
	ChannelId   string           // 通道id
	ChannelType enum.ChannelType // 通道类型
	clientId    uint             // 客户端id
	conn        *websocket.Conn  // 前端websocket连接
	isClosed    bool             // 标记是否已经关闭
	serverPort  string           // 服务端端口
	tcpConn     net.Conn         // tcp连接
	clientAddr  string
}

func (c Service) NewChannel(id uint, channelType enum.ChannelType, conn *websocket.Conn, serverPort string, clientAddr string) (err error) {
	_, err = c.getClientWebsocketConn(id)
	if err != nil {
		return err
	}

	if channelType == enum.Tcp {
		// 服务端需要监听新端口
		go c.listen(id, serverPort, clientAddr)
	}

	if channelType == enum.Pty {
		channel, err := c.createChannel(id, enum.Pty, "")
		if err != nil {
			return err
		}
		channel.conn = conn
		go channel.ptyRead(c)
		go channel.ptyWrite(c)
	}

	return nil
}

func (c Service) createChannel(id uint, channelType enum.ChannelType, clientAddr string) (channel *Channel, err error) {
	channelId := utils.RandString(16)
	channel = &Channel{
		ChannelId:   channelId,
		ChannelType: channelType,
		clientId:    id,
		isClosed:    false,
		clientAddr:  clientAddr,
	}
	c.channels[channelId] = channel
	c.messageMq[channel.ChannelId] = make(chan []byte, 32)

	// 通知客户端打开对应通道
	_, err = c.SendCommand(id, enum.MessageTypeChannel, &request.ChannelRequest{
		Action:      "open",
		ChannelId:   channel.ChannelId,
		ChannelType: channelType,
		ChannelData: nil,
		Addr:        clientAddr,
	})
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (c Service) listen(id uint, serverPort string, clientAddr string) {
	log.Info("server start listening on port "+serverPort, nil)
	listener, err := net.Listen("tcp", ":"+serverPort)
	if err != nil {
		log.Error("listen error: "+err.Error(), nil)
		return
	}
	defer listener.Close()

	// 监听连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error("Error accepting connection:"+err.Error(), nil)
			continue
		}
		fmt.Println("New connection from:", conn.RemoteAddr())

		channel, err := c.createChannel(id, enum.Tcp, clientAddr)
		if err != nil {
			log.Error("Error creating channel:"+err.Error(), nil)
			continue
		}

		channel.tcpConn = conn

		go channel.tcpRead(c)
		go channel.tcpWrite(c)
	}
}

func (channel *Channel) tcpRead(c Service) {
	defer channel.close()
	buffer := make([]byte, 1024)
	for {
		n, err := channel.tcpConn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from TCP connection:", err)
			return
		}
		ws_data := request.ChannelRequest{
			Action:      "write",
			ChannelId:   channel.ChannelId,
			ChannelData: buffer[:n],
		}
		_, err = c.clientWebsocketWrite(channel.clientId, websocket.TextMessage, enum.MessageTypeChannel, ws_data, "")
		if err != nil {
			fmt.Println("Error writing to WebSocket:", err)
			return
		}
	}
}

func (channel *Channel) tcpWrite(c Service) error {
	for {
		data := <-c.messageMq[channel.ChannelId]
		if _, err := channel.tcpConn.Write(data); err != nil {
			log.Error("Error writing to TCP:"+err.Error(), nil)
			return err
		}
	}
	return nil
}

func (channel *Channel) ptyWrite(c Service) error {
	for {
		data := <-c.messageMq[channel.ChannelId]
		if err := channel.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
			log.Error("Error writing to TCP:"+err.Error(), nil)
			return err
		}
	}
	return nil
}

func (channel *Channel) ptyRead(c Service) {
	defer channel.close()
	for {
		msgType, data, err := channel.conn.ReadMessage()
		if err != nil {
			log.Error("channel read goroutine ReadMessage error: "+err.Error(), nil)
			return
		}

		_, err = c.clientWebsocketWrite(channel.clientId, msgType, enum.MessageTypeChannel, &request.ChannelRequest{
			Action:      "write",
			ChannelId:   channel.ChannelId,
			ChannelData: data,
		}, "")
		if err != nil {
			return
		}
	}
}

func (channel *Channel) close() error {
	if channel.isClosed {
		return nil
	}

	if channel.conn != nil {
		channel.conn.Close()
	}

	if channel.tcpConn != nil {
		channel.tcpConn.Close()
	}

	channel.isClosed = true

	return nil
}
