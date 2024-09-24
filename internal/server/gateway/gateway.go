package gateway

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"noah/internal/server/enum"
	"noah/internal/server/middleware/log"
	"noah/internal/server/request"
	"noah/internal/server/utils"
	"strconv"
	"sync"
)

type Gateway struct {
	mu        *sync.Mutex
	clients   map[uint]*websocket.Conn //客户端命令执行websocket
	messageMq map[string]chan request.Message
	messageId uint64 //消息id
}

func NewGateway() *Gateway {
	g := &Gateway{
		mu:        &sync.Mutex{},
		clients:   make(map[uint]*websocket.Conn),
		messageMq: make(map[string]chan request.Message),
		messageId: 0,
	}

	//go func() {
	//	for {
	//		log.Info("gateway health check", map[string]interface{}{
	//			"clients":   g.clients,
	//			"messageMq": g.messageMq,
	//		})
	//
	//		time.Sleep(time.Second * 10)
	//	}
	//}()
	return g
}

// NewClientWebsocketConn 新增ws连接
func (g Gateway) NewClientWebsocketConn(clientId uint, connection *websocket.Conn) error {
	// 验证连接是否有效
	if connection == nil || connection.RemoteAddr() == nil {
		return errors.New("invalid connection")
	}
	g.mu.Lock()
	defer g.mu.Unlock()

	// 更新或添加新连接
	g.clients[clientId] = connection

	// 启动一个goroutine用于读取客户端的websocket消息
	go g.clientWebsocketRead(clientId)
	return nil
}

// getClientWebsocketConn 获取客户端websocket连接
func (g Gateway) getClientWebsocketConn(clientID uint) (*websocket.Conn, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	conn, found := g.clients[clientID]
	if !found {
		return nil, errors.New("client not found")
	}
	return conn, nil
}

// Exit 关闭连接
func (g Gateway) exit(clientId uint) error {
	if err := g.removeClientWebsocketConnection(clientId); err != nil {
		return err
	}
	return nil
}

// removeClientWebsocketConnection 删除客户端websocket连接
func (g Gateway) removeClientWebsocketConnection(clientID uint) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if conn, found := g.clients[clientID]; !found {
	} else {
		err := conn.Close()
		if err != nil {
		}
		delete(g.clients, clientID)
	}

	return nil
}

// ClientWebsocketWrite 向发送消息
func (g Gateway) ClientWebsocketWrite(clientId uint, messageId string, messageType enum.MessageType, data any) (string, error) {
	client, err := g.getClientWebsocketConn(clientId)
	if err != nil {
		return "", err
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
			return "", err
		}
	}

	if messageId == "" {
		messageId = strconv.FormatUint(g.messageId+1, 10)
	}

	req := request.Message{
		MessageId:   messageId,
		MessageType: messageType,
		Data:        messageByteData,
		Error:       "",
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	//todo 消息加密

	err = client.WriteMessage(websocket.TextMessage, body)
	if err != nil {
		log.Error("write client websocket message error, websocket exit", map[string]interface{}{"clientId": clientId, "error": err})
		err := g.exit(clientId)
		if err != nil {
			log.Error("remove client websocket connection error", map[string]interface{}{"clientId": clientId, "error": err})
			return "", err
		}
		return "", err
	}

	if _, ok := g.messageMq[messageId]; !ok {
		g.messageMq[messageId] = make(chan request.Message, 32)
	}

	return messageId, nil
}

// clientWebsocketRead 读取客户端的websocket消息
func (g Gateway) clientWebsocketRead(clientId uint) {
	defer func() {
		log.Info("client websocket read fail, websocket exit", map[string]interface{}{"clientId": clientId})
		err := g.exit(clientId)
		if err != nil {
			log.Error("remove client websocket connection error", map[string]interface{}{"clientId": clientId, "error": err})
			return
		}
	}()
	client, err := g.getClientWebsocketConn(clientId)
	if err != nil {
		return
	}
	for {
		_, wsMessage, err := client.ReadMessage()
		if err != nil {
			log.Error("read message error", map[string]interface{}{"clientId": clientId, "error": err})
			return
		}

		//todo 消息加密

		var message request.Message
		if err := json.Unmarshal(wsMessage, &message); err != nil {
			log.Error("parse message error", map[string]interface{}{"clientId": clientId, "error": err})
			continue
		}

		//fixme 消息转发
		if _, ok := g.messageMq[message.MessageId]; ok {
			g.messageMq[message.MessageId] <- message
		}
	}
}

func (g Gateway) ClientWebsocketRead(messageId string) (request.Message, error) {
	if _, ok := g.messageMq[messageId]; ok {
		return <-g.messageMq[messageId], nil
	}
	return request.Message{}, errors.New("message not found")
}

func (g Gateway) CloseMessageMq(messageId string) {
	if _, ok := g.messageMq[messageId]; ok {
		close(g.messageMq[messageId])
		delete(g.messageMq, messageId)
	}
}

// SendCommand 执行命令
func (g Gateway) SendCommand(clientId uint, messageType enum.MessageType, data any) (string, error) {
	msgId, err := g.ClientWebsocketWrite(clientId, "", messageType, data)
	if err != nil {
		return "", err
	}

	// 获取命令执行结果
	result, err := g.ClientWebsocketRead(msgId)
	defer g.CloseMessageMq(msgId)
	if err != nil {
		return "", err
	}

	if result.Error != "" {
		return "", errors.New(result.Error)
	}

	return utils.ByteToString(result.Data), nil
}
