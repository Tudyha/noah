package app

import (
	"context"
	"log"
	"net"
	"noah/client/app/config"
	pkgApp "noah/pkg/app"
	"noah/pkg/conn"
	"noah/pkg/packet"
	"sync"
	"time"
)

type Client struct {
	cfg *config.Config

	mux       sync.Mutex
	connected bool

	conn *conn.Conn

	closeSignal chan struct{}
}

func NewClient(cfg *config.Config) pkgApp.Server {
	return &Client{
		cfg:         cfg,
		closeSignal: make(chan struct{}),
	}
}

func (c *Client) Start(ctx context.Context) error {
	for {
		select {
		case <-c.closeSignal:
			return nil
		default:
			c.mux.Lock()
			defer c.mux.Unlock()
			if c.connected {
				continue
			}
			if err := c.connect(); err != nil {
				log.Println("连接失败:", err)
				time.Sleep(10 * time.Second)
			}
			c.connected = true
		}
	}
}

func (c *Client) connect() error {
	dail := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	netConn, err := dail.Dial("tcp", c.cfg.ServerAddr)
	if err != nil {
		return err
	}
	c.conn = conn.NewConn(netConn)
	go c.handleConn()
	return nil
}

func (c *Client) handleConn() {
	defer func() {
		c.mux.Lock()
		defer c.mux.Unlock()
		c.connected = false
		c.conn.Close()
	}()

	// 发送注册包
	msg := packet.Login{
		AppId:    c.cfg.AppId,
		Sign:     "123456",
		DeviceId: "1",
	}
	if err := c.conn.WriteMessage(packet.MessageType_Login, &msg); err != nil {
		return
	}
	log.Println("注册包发送成功")
	for {

	}
}

func (c *Client) Stop(ctx context.Context) error {
	if c.conn != nil {
		c.conn.Close()
	}
	close(c.closeSignal)
	return nil
}

func (c *Client) String() string {
	return "client"
}
