package channel

import (
	"encoding/json"
	"errors"
	"github.com/duke-git/lancet/v2/slice"
	"net"
	"noah/internal/server/config"
	"noah/internal/server/enum"
	"noah/internal/server/middleware/log"
	"noah/internal/server/request"
	"noah/internal/server/utils"
	"sync"
	"time"

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
}

func NewChannelService() *Service {
	return &Service{
		mu:            &sync.Mutex{},
		clients:       make(map[uint]*websocket.Conn),
		messageResult: make(map[uint64]chan request.Message),
		channels:      make(map[string]*Channel),
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

func (c Service) clientWebsocketWrite(id uint, websocketMessageType int, messageType enum.MessageType, data any, errMsg string) (messageId uint64, err error) {
	client, err := c.getClientWebsocketConn(id)
	if err != nil {
		return 0, err
	}

	var d []byte
	switch data.(type) {
	case []byte:
		d = data.([]byte)
	case string:
		d = []byte(data.(string))
	default:
		d, err = json.Marshal(data)
		if err != nil {
			return 0, err
		}
	}

	messageId = messageId + 1

	body, err := json.Marshal(request.Message{
		MessageId:   messageId,
		MessageType: messageType,
		Data:        d,
		Error:       errMsg,
	})
	if err != nil {
		return 0, err
	}

	err = client.WriteMessage(websocketMessageType, body)
	if err != nil {
		return 0, err
	}

	return messageId, nil
}

func (c Service) clientWebsocketRead(id uint) {
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
			continue
		}
		if message.Error != "" {
			log.Warn("receive message error", map[string]interface{}{"clientId": id, "error": err})
			continue
		}
		c.handleMessage(id, message)
	}
}

func (c Service) handleMessage(id uint, message request.Message) {
	switch message.MessageType {
	case enum.MessageTypePty:
		var ptyRequest request.PtyRequest
		err := json.Unmarshal(message.Data, &ptyRequest)
		if err != nil {
			log.Error("parse pty request error", map[string]interface{}{"clientId": id, "error": err})
			return
		}
		switch ptyRequest.Action {
		case "write":
			c.write(ptyRequest.ChannelId, ptyRequest.ChannelData)
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
	if messageType == enum.MessageTypePty {

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
	//closeSignal chan struct{}    // 用于关闭连接的信号
	isClosed   bool       // 标记是否已经关闭
	mu         sync.Mutex // 用于保护对channelId的访问
	serverPort string     // 服务端端口
}

func (c Service) NewChannel(id uint, channelType enum.ChannelType, conn *websocket.Conn, serverPort string) (channel *Channel, err error) {
	_, err = c.getClientWebsocketConn(id)
	if err != nil {
		return nil, err
	}

	channelId := utils.RandString(16)
	channel = &Channel{
		ChannelId:   channelId,
		ChannelType: channelType,
		conn:        conn,
		isClosed:    false,
		//closeSignal: make(chan struct{}),
		serverPort: serverPort,
	}

	if channelType == enum.Tcp {
		// 服务端需要监听新端口
		go channel.listen(c, id)
	}

	if channelType == enum.Pty {
		// pty通道需要通知客户端打开pty TODO need result to check
		_, err = c.SendCommand(id, enum.MessageTypePty, &request.PtyRequest{
			Action:      "open",
			ChannelId:   channel.ChannelId,
			ChannelData: nil,
		})
		if err != nil {
			return nil, err
		}

		go channel.read(c, id)
	}

	c.channels[channelId] = channel
	return channel, nil
}

func (channel *Channel) listen(c Service, id uint) {
	port := channel.serverPort
	log.Info("server start listening on port "+port, nil)
	listener, err := net.Listen("tcp", ":"+port)
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

		go channel.tcpRead(conn, c, id)
		go channel.tcpWrite(conn, c, id)
	}
}

func (channel *Channel) tcpRead(tcpConn net.Conn, c Service, id uint) {
	// 从 TCP 读取数据并写入 WebSocket
	buffer := make([]byte, 1024)
	for {
		n, err := tcpConn.Read(buffer)
		if err != nil {
			log.Error("Error reading from TCP:"+err.Error(), nil)
			break
		}
		c.clientWebsocketWrite(id, websocket.TextMessage, enum.MessageTypeChannel, buffer[:n], "")
	}
}

func (channel *Channel) tcpWrite(tcpConn net.Conn, c Service, id uint) {
	//for {
	//	_, message, err := clientWebsocketConn.ReadMessage()
	//	if err != nil {
	//		log.Error("Error reading from WebSocket:"+err.Error(), nil)
	//		break
	//	}
	//
	//	if _, err = tcpConn.Write(message); err != nil {
	//		log.Error("Error writing to TCP:"+err.Error(), nil)
	//		break
	//	}
	//}
}

func (c Service) write(channelId string, data []byte) {
	channel, ok := c.channels[channelId]
	if !ok {
		return
	}

	channel.conn.SetWriteDeadline(time.Now().Add(config.MessageWait))
	if err := channel.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		log.Error("channel write goroutine WriteMessage error: "+err.Error(), nil)
		return
	}

}

func (channel *Channel) read(c Service, id uint) {
	for {
		msgType, data, err := channel.conn.ReadMessage()
		if err != nil {
			log.Error("channel read goroutine ReadMessage error: "+err.Error(), nil)
			return
		}

		_, err = c.clientWebsocketWrite(id, msgType, enum.MessageTypePty, &request.PtyRequest{
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
	channel.mu.Lock()
	defer channel.mu.Unlock()

	if channel.isClosed {
		return nil
	}

	if channel.conn != nil {
		channel.conn.Close()
	}

	channel.isClosed = true

	return nil
}
