package app

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"noah/client/app/handler"
	pkgApp "noah/pkg/app"
	"noah/pkg/config"
	"noah/pkg/conn"
	"noah/pkg/packet"
	"noah/pkg/utils"
	"sync/atomic"
	"time"

	"github.com/xtaci/smux/v2"
)

type Client struct {
	cfg *config.ClientConfig

	connected atomic.Bool

	conn       *conn.Conn
	session    *smux.Session
	pingStream *smux.Stream

	closeSignal chan struct{}

	infoHandler *handler.InfoHandler

	messageHandlers map[packet.MessageType]conn.MessageHandler
}

func NewClient(cfg *config.ClientConfig) pkgApp.Server {
	if cfg.ReconnectInterval == 0 {
		cfg.ReconnectInterval = 10
	}
	if cfg.DailTimeout == 0 {
		cfg.DailTimeout = 10
	}
	if cfg.HeartbeatInterval == 0 {
		cfg.HeartbeatInterval = 30
	}
	c := &Client{
		cfg:         cfg,
		closeSignal: make(chan struct{}),
		infoHandler: &handler.InfoHandler{},

		messageHandlers: make(map[packet.MessageType]conn.MessageHandler),
	}

	logoutHandler := handler.NewLogoutHandler()
	c.messageHandlers[logoutHandler.MessageType()] = logoutHandler

	return c
}

func (c *Client) Start(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(c.cfg.ReconnectInterval) * time.Second)
	defer ticker.Stop()

	// 主循环：监听关闭信号并维持客户端连接
	for {
		select {
		case <-c.closeSignal:
			// 接收到关闭信号，停止客户端
			return nil
		case <-ticker.C:
			if !c.connected.Load() {
				if err := c.connect(); err != nil {
					log.Println("连接失败:", err)
				}
			}
		}
	}
}

func (c *Client) connect() error {
	dail := net.Dialer{
		Timeout: time.Duration(c.cfg.DailTimeout) * time.Second,
	}

	netConn, err := dail.Dial("tcp", c.cfg.ServerAddr)
	if err != nil {
		return err
	}

	c.conn = conn.NewConn(netConn)
	c.connected.Store(true)

	go c.handleConn()
	return nil
}

func (c *Client) handleConn() {
	defer func() {
		c.connected.Store(false)
		if c.conn != nil {
			c.conn.Close()
		}
		if c.session != nil {
			c.session.Close()
		}
	}()

	// 发送登录包
	loginReq := &packet.Login{
		AppId:    c.cfg.AppId,
		Sign:     utils.Sign(c.cfg.AppId, c.cfg.AppSecret),
		DeviceId: utils.GetMacAddress(),
	}
	loginReq.ClientInfo = c.infoHandler.GetInfo()
	if _, err := c.conn.WriteProtoMessage(packet.MessageType_Login, loginReq); err != nil {
		return
	}
	p, err := c.conn.ReadOnce()
	if err != nil {
		return
	}
	if p.MessageType != packet.MessageType_LoginAck {
		return
	}

	log.Println("认证通过")

	log.Println("创建smux session")
	c.session, err = smux.Server(c.conn.GetConn(), smux.DefaultConfig())
	if err != nil {
		log.Println("创建smux session失败:", err)
		return
	}
	c.conn.Stop()
	c.conn = nil

	log.Println("创建ping任务")
	if c.pingStream, err = c.session.OpenStream(); err != nil {
		log.Println("创建ping stream失败:", err)
	} else {
		go c.ping()
	}

	log.Println("监听新连接")
	for {
		stream, err := c.session.AcceptStream()
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				return
			}
			continue
		}
		go c.handleStream(stream)
	}
}

func (c *Client) handleStream(conn net.Conn) {
	log.Println("handle new stream")
}

func (c *Client) ping() {
	if c.pingStream == nil {
		return
	}
	conn := conn.NewConn(c.pingStream)
	defer conn.Close()

	ticker := time.NewTicker(time.Duration(c.cfg.HeartbeatInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.closeSignal:
			log.Println("client close, stop ping")
			return
		case <-ticker.C:
			if !c.connected.Load() {
				log.Println("client not connected, ping fail")
				continue
			}
			data := c.infoHandler.GetSystemStat()
			if _, err := conn.WriteProtoMessage(packet.MessageType_Ping, data); err != nil {
				log.Println("心跳包发送失败:", err)
				return
			}
			log.Println("心跳包发送成功")
		}
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
