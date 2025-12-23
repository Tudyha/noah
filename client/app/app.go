package app

import (
	"context"
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
	tunnelHandler := handler.NewTunnelHandler()
	c.messageHandlers[logoutHandler.MessageType()] = logoutHandler
	c.messageHandlers[tunnelHandler.MessageType()] = tunnelHandler

	return c
}

func (c *Client) Start(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(c.cfg.ReconnectInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.closeSignal:
			// 接收到关闭信号，停止客户端
			return nil
		case <-ticker.C:
			// 定时重连
			if !c.connected.Load() {
				if err := c.connect(); err != nil {
					log.Println("连接失败:", err)
				}
			}
		}
	}
}

// 连接到服务端
func (c *Client) connect() error {
	dail := net.Dialer{
		Timeout: time.Duration(c.cfg.DailTimeout) * time.Second,
	}

	netConn, err := dail.Dial("tcp", c.cfg.ServerAddr)
	if err != nil {
		return err
	}

	c.connected.Store(true)

	go c.handleConn(netConn)
	return nil
}

// 处理连接
func (c *Client) handleConn(netConn net.Conn) {
	conn := conn.NewConn(netConn)
	defer func() {
		c.connected.Store(false)
		if conn != nil {
			conn.Close()
		}
		if c.session != nil {
			c.session.Close()
		}
		if c.pingStream != nil {
			c.pingStream.Close()
		}
	}()

	// 发送鉴权包
	loginReq := &packet.Login{
		AppId:    c.cfg.AppId,
		Sign:     utils.Sign(c.cfg.AppId, c.cfg.AppSecret),
		DeviceId: utils.GetMacAddress(),
	}
	loginReq.ClientInfo = c.infoHandler.GetInfo()
	if err := conn.WriteProtoMessage(packet.MessageType_Login, loginReq); err != nil {
		return
	}

	// 读取鉴权包ack
	p, err := conn.ReadMessage()
	if err != nil {
		return
	}
	if p.MessageType != packet.MessageType_LoginAck {
		return
	}

	log.Println("认证通过")

	conn.Release()
	conn = nil

	log.Println("创建smux session")
	c.session, err = smux.Client(netConn, smux.DefaultConfig())
	if err != nil {
		log.Println("创建smux session失败:", err)
		return
	}

	log.Println("smux session创建成功，创建ping stream")
	if c.pingStream, err = c.session.OpenStream(); err != nil {
		log.Println("创建ping stream失败:", err)
	} else {
		log.Println("开始定时发送心跳包")
		go c.ping()
	}

	for {
		stream, err := c.session.AcceptStream()
		if err != nil {
			return
		}
		go c.handleStream(stream)
	}
}

// 处理stream
func (c *Client) handleStream(netConn net.Conn) {
	co := conn.NewConn(netConn)
	for {
		p, err := co.ReadMessage()
		if err != nil {
			co.Close()
			return
		}
		// 创建上下文，用于传递数据
		ctx := conn.NewConnContext(co, p)
		h := c.messageHandlers[p.MessageType]
		if h == nil {
			log.Println("消息处理函数未注册", "messageType", p.MessageType)
			ctx.Release()
			continue
		}
		if err = h.Handle(ctx); err != nil {
			log.Println("处理消息失败", "err", err)
		}

		if ctx.IsHijacked() {
			log.Println("stream底层连接已被劫持，不再处理消息")
			ctx.Release()
			return
		}

		ctx.Release()
	}
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
			if err := conn.WriteProtoMessage(packet.MessageType_Ping, data); err != nil {
				log.Println("心跳包发送失败:", err)
				return
			}
			log.Println("心跳包发送成功")
		}
	}
}

func (c *Client) Stop(ctx context.Context) error {
	log.Println("client stop...")
	if c.session != nil {
		c.session.Close()
	}
	close(c.closeSignal)
	return nil
}

func (c *Client) String() string {
	return "client"
}
