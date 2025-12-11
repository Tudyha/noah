package app

import (
	"context"
	"log"
	"net"
	"noah/client/app/config"
	pkgApp "noah/pkg/app"
	"noah/pkg/conn"
	"noah/pkg/packet"
	"noah/pkg/utils"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Client struct {
	cfg *config.Config

	mux       sync.Mutex
	connected atomic.Bool

	conn *conn.Conn

	closeSignal chan struct{}
}

func NewClient(cfg *config.Config) pkgApp.Server {
	if cfg.ReconnectInterval == 0 {
		cfg.ReconnectInterval = 10
	}
	return &Client{
		cfg:         cfg,
		closeSignal: make(chan struct{}),
	}
}

func (c *Client) Start(ctx context.Context) error {
	log.Println("client start...")
	for {
		select {
		case <-c.closeSignal:
			return nil
		default:
			if c.connected.Load() {
				continue
			}
			if err := c.connect(); err != nil {
				log.Println("连接失败:", err)
				time.Sleep(time.Duration(c.cfg.ReconnectInterval) * time.Second)
				continue
			}
			c.connected.Store(true)
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
		c.connected.Store(false)
		c.conn.Close()
	}()

	// 发送登录包
	body, err := anypb.New(&packet.Login{
		AppId:    c.cfg.AppId,
		Sign:     utils.Sign(c.cfg.AppId, "123456"),
		DeviceId: "1",
	})
	if err != nil {
		return
	}

	data, err := proto.Marshal(&packet.Message{
		Body: body,
	})
	if err != nil {
		return
	}
	if err := c.conn.WriteProtoMessage(packet.MessageType_Login, data); err != nil {
		return
	}
	log.Println("注册包发送成功")
	for {

	}
}

func (c *Client) Stop(ctx context.Context) error {
	log.Println("client stop...")
	if c.conn != nil {
		c.conn.Close()
	}
	close(c.closeSignal)
	return nil
}

func (c *Client) String() string {
	return "client"
}
