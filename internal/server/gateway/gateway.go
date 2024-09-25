package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"noah/internal/server/enum"
	"noah/internal/server/middleware/log"
	"noah/internal/server/request"
	"noah/internal/server/utils"
	"strconv"
	"sync"
	"time"
)

type Gateway struct {
	mu                 *sync.Mutex
	clients            map[uint]*websocket.Conn //客户端命令执行websocket
	messageId          uint64                   //消息id
	messageSubscribers map[uint][]messageSubscriber
}

type messageSubscriber struct {
	messageId string
	f         func(message request.Message)
}

func NewGateway() *Gateway {
	g := &Gateway{
		mu:                 &sync.Mutex{},
		clients:            make(map[uint]*websocket.Conn),
		messageId:          0,
		messageSubscribers: make(map[uint][]messageSubscriber),
	}

	//go func() {
	//	for {
	//		log.Info("gateway health check", map[string]interface{}{
	//			"clients":            g.clients,
	//			"messageSubscribers": g.messageSubscribers[1],
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
	g.messageSubscribers[clientId] = make([]messageSubscriber, 0)

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

	if conn, found := g.clients[clientID]; found {
		err := conn.Close()
		if err != nil {
		}
		delete(g.clients, clientID)
	}

	if _, found := g.messageSubscribers[clientID]; found {
		delete(g.messageSubscribers, clientID)
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
		g.messageId++
		messageId = strconv.FormatUint(g.messageId, 10)
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

		for _, subscriber := range g.messageSubscribers[clientId] {
			subscriber.f(message)
		}
	}
}

func (g Gateway) SubscribeMessage(clientId uint, messageId string, f func(message request.Message)) {
	g.messageSubscribers[clientId] = append(g.messageSubscribers[clientId], messageSubscriber{
		messageId: messageId,
		f:         f,
	})
}

func (g Gateway) UnSubscribeMessage(clientId uint, messageId string) {
	for i, subscriber := range g.messageSubscribers[clientId] {
		if subscriber.messageId == messageId {
			g.messageSubscribers[clientId] = append(g.messageSubscribers[clientId][:i], g.messageSubscribers[clientId][i+1:]...)
			break
		}
	}
}

// SendCommand 执行命令
func (g Gateway) SendCommand(clientId uint, messageType enum.MessageType, data any, needResult bool) (string, error) {
	msgId, err := g.ClientWebsocketWrite(clientId, "", messageType, data)
	if err != nil {
		return "", err
	}
	if !needResult {
		return "", nil
	}

	// 创建一个带超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result request.Message
	done := make(chan struct{})

	// 使用 context 来控制订阅消息的 goroutine 的生命周期
	go func() {
		g.SubscribeMessage(clientId, msgId, func(message request.Message) {
			result = message
			close(done)
		})
		defer g.UnSubscribeMessage(clientId, msgId)

		select {
		case <-ctx.Done():
		}
	}()

	select {
	case <-done:

	case <-ctx.Done():
		return "", errors.New("timeout")
	}

	if result.Error != "" {
		return "", errors.New(result.Error)
	}

	return utils.ByteToString(result.Data), nil
}
