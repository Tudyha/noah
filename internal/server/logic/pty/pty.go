package pty

import (
	"errors"
	"noah/internal/server/service"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ptyService struct {
	mu         *sync.Mutex
	ptyClients map[string]*websocket.Conn //客户端pty websocket
}

func init() {
	service.RegisterPtyService(&ptyService{
		mu:         &sync.Mutex{},
		ptyClients: make(map[string]*websocket.Conn),
	})
}

type ptyClient struct {
	conn        *websocket.Conn
	channelId   string
	closeSignal chan struct{}
	isClosed    bool       // 标记是否已经关闭
	mu          sync.Mutex // 用于保护对channelId的访问
}

const (
	// Time allowed to write or read a message.
	messageWait = 10 * time.Second
)

var (
	ErrClientConnectionNotFound = errors.New("no active client connection found")
)

func (c *ptyService) PtyRead(channelId string) (messageType int, data []byte, err error) {
	client, found := c.getPtyConnection(channelId)
	if !found {
		return 0, nil, ErrClientConnectionNotFound
	}

	msgType, connReader, err := client.ReadMessage()
	if err != nil {
		return 0, nil, err
	}

	return msgType, connReader, nil
}

func (c *ptyService) PtyWrite(channelId string, messageType int, data []byte) error {
	client, _ := c.getPtyConnection(channelId)

	err := client.WriteMessage(messageType, data)
	if err != nil {
		return err
	}
	return nil
}

// ClosePtyConnection 关闭pty连接
func (c *ptyService) ClosePtyConnection(channelId string) error {
	client, _ := c.getPtyConnection(channelId)
	if client != nil {
		if err := client.Close(); err != nil {
			return err
		}
	}

	if err := c.removePtyConnection(channelId); err != nil {
		return err
	}

	return nil
}

func (c *ptyService) AddPtyConnection(channelId string, connection *websocket.Conn) error {
	c.mu.Lock()
	c.ptyClients[channelId] = connection
	c.mu.Unlock()
	return nil
}

func (c *ptyService) getPtyConnection(channelId string) (*websocket.Conn, bool) {
	c.mu.Lock()
	conn, found := c.ptyClients[channelId]
	c.mu.Unlock()
	return conn, found
}

func (c *ptyService) removePtyConnection(channelId string) error {
	c.mu.Lock()
	delete(c.ptyClients, channelId)
	c.mu.Unlock()
	return nil
}

func (p *ptyService) NewPtyClient(channelId string, conn *websocket.Conn) error {
	if conn == nil {
		return errors.New("connection is nil")
	}

	ptyClient := &ptyClient{
		conn:        conn,
		channelId:   channelId,
		closeSignal: make(chan struct{}),
		isClosed:    false,
	}

	go ptyClient.read()
	go ptyClient.write()
	return nil
}

// 往前端发送消息
func (c *ptyClient) write() {
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
			//TODO ClientService问题
			msgType, data, err := service.GetPtyService().PtyRead(c.channelId)
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
func (c *ptyClient) read() {
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
			service.GetPtyService().PtyWrite(c.channelId, msgType, data)
			c.mu.Unlock()
		}
	}
}

func (c *ptyClient) close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isClosed {
		return nil
	}

	if c.conn != nil {
		c.conn.Close()
	}

	//断开与client的websocket连接
	service.GetPtyService().ClosePtyConnection(c.channelId)

	close(c.closeSignal)
	c.isClosed = true

	return nil
}
