package pty

import (
	"errors"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ptyService struct {
	mu         *sync.Mutex
	ptyClients map[string]*websocket.Conn //客户端pty websocket
}

func NewPtyService() *ptyService {
	return &ptyService{
		mu:         &sync.Mutex{},
		ptyClients: make(map[string]*websocket.Conn),
	}
}

type ptyChannel struct {
	channelId   string          // 通道id
	conn        *websocket.Conn // 前端websocket连接
	closeSignal chan struct{}   // 用于关闭连接的信号
	isClosed    bool            // 标记是否已经关闭
	mu          sync.Mutex      // 用于保护对channelId的访问
	ptyService  *ptyService     // ptyService，用于操作ptyClient
}

const (
	// Time allowed to write or read a message.
	messageWait = 10 * time.Second
)

var (
	ErrClientConnectionNotFound = errors.New("no active client connection found")
)

func (c ptyService) NewPtyClient(channelId string, connection *websocket.Conn) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ptyClients[channelId] = connection
	return nil
}

func (c ptyService) PtyClientRead(channelId string) (messageType int, data []byte, err error) {
	client, found := c.getPtyClientConnection(channelId)
	if !found {
		return 0, nil, ErrClientConnectionNotFound
	}

	msgType, connReader, err := client.ReadMessage()
	if err != nil {
		return 0, nil, err
	}

	return msgType, connReader, nil
}

func (c ptyService) PtyClientWrite(channelId string, messageType int, data []byte) error {
	client, _ := c.getPtyClientConnection(channelId)

	err := client.WriteMessage(messageType, data)
	if err != nil {
		return err
	}
	return nil
}

// ClosePtyClient 关闭pty连接
func (c ptyService) ClosePtyClient(channelId string) error {
	client, _ := c.getPtyClientConnection(channelId)
	if client != nil {
		if err := client.Close(); err != nil {
			return err
		}
	}

	if err := c.removePtyClientConnection(channelId); err != nil {
		return err
	}

	return nil
}

func (c ptyService) getPtyClientConnection(channelId string) (*websocket.Conn, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	conn, found := c.ptyClients[channelId]
	return conn, found
}

func (c ptyService) removePtyClientConnection(channelId string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.ptyClients, channelId)
	return nil
}

// NewPtyChannel 新建pty通道
func (c ptyService) NewPtyChannel(channelId string, conn *websocket.Conn) error {
	if conn == nil {
		return errors.New("connection is nil")
	}

	channel := &ptyChannel{
		conn:        conn,
		channelId:   channelId,
		closeSignal: make(chan struct{}),
		isClosed:    false,
		ptyService:  &c,
	}

	go channel.read()
	go channel.write()
	return nil
}

// 往前端发送消息
func (c *ptyChannel) write() {
	defer func() {
		if err := c.close(); err != nil {
			//log.Fatalf("ptyClient.close: %s", err)
		}
	}()

	for {
		select {
		case <-c.closeSignal:
			return
		default:
			time.Sleep(10 * time.Millisecond)
			// 读取client发来的消息
			msgType, data, err := c.ptyService.PtyClientRead(c.channelId)
			if err != nil {
				//log.Fatalf("ptyClient.write: %s", err)
				return
			}

			//转发到前端
			c.conn.SetWriteDeadline(time.Now().Add(messageWait))
			if err := c.conn.WriteMessage(msgType, data); err != nil {
				//log.Fatalf("conn.WriteMessage: %s", err)
				return
			}
		}
	}
}

// 从前端读取消息
func (c *ptyChannel) read() {
	defer func() {
		if err := c.close(); err != nil {
			//log.Fatalf("ptyClient.close: %s", err)
		}
	}()

	for {
		select {
		case <-c.closeSignal:
			return
		default:
			// 从前端读取消息
			msgType, data, err := c.conn.ReadMessage()
			if err != nil {
				//log.Fatalf("conn.ReadMessage: %s", err)
				return
			}

			// 转发到 client
			c.mu.Lock()
			c.ptyService.PtyClientWrite(c.channelId, msgType, data)
			c.mu.Unlock()
		}
	}
}

func (c *ptyChannel) close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isClosed {
		return nil
	}

	if c.conn != nil {
		c.conn.Close()
	}

	//断开与client的websocket连接
	c.ptyService.ClosePtyClient(c.channelId)

	close(c.closeSignal)
	c.isClosed = true

	return nil
}
