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

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Client struct {
	cfg *config.ClientConfig

	connected atomic.Bool

	conn *conn.Conn

	closeSignal chan struct{}

	infoHandler *handler.InfoHandler
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
	return &Client{
		cfg:         cfg,
		closeSignal: make(chan struct{}),
		infoHandler: &handler.InfoHandler{},
	}
}

// Start 启动客户端并维持连接状态
// 参数:
//
//	ctx: 上下文对象，用于控制协程生命周期
//
// 返回值:
//
//	error: 启动过程中发生的错误，正常关闭时返回nil
func (c *Client) Start(ctx context.Context) error {
	log.Println("client start...")

	// 定时发送心跳包
	go c.ping()

	// 主循环：监听关闭信号并维持客户端连接
	for {
		select {
		case <-c.closeSignal:
			// 接收到关闭信号，停止客户端
			log.Println("client stop...")
			return nil
		default:
			// 检查客户端是否已连接，如果已连接则跳过本次循环
			if c.connected.Load() {
				continue
			}

			// 尝试建立连接，如果失败则等待重连间隔后重试
			if err := c.connect(); err != nil {
				log.Println("连接失败:", err)
				time.Sleep(time.Duration(c.cfg.ReconnectInterval) * time.Second)
				continue
			}
		}
	}
}

// connect 建立与服务器的TCP连接并初始化连接处理器
// 该函数会创建一个带有超时设置的网络拨号器，尝试连接到配置的服务器地址，
// 如果连接成功，则初始化连接对象并启动连接处理协程
//
// 返回值:
//
//	error - 连接过程中的错误信息，成功时返回nil
func (c *Client) connect() error {
	// 创建带超时设置的拨号器
	dail := net.Dialer{
		Timeout: time.Duration(c.cfg.DailTimeout) * time.Second,
	}

	// 尝试建立TCP连接
	netConn, err := dail.Dial("tcp", c.cfg.ServerAddr)
	if err != nil {
		return err
	}

	// 初始化连接并标记为已连接状态
	c.conn = conn.NewConn(netConn)
	c.connected.Store(true)

	// 启动连接处理协程
	go c.handleConn()
	return nil
}

func (c *Client) handleConn() {
	defer func() {
		c.connected.Store(false)
		c.conn.Close()
	}()

	// 发送登录包
	loginReq := &packet.Login{
		AppId:    c.cfg.AppId,
		Sign:     utils.Sign(c.cfg.AppId, c.cfg.AppSecret),
		DeviceId: "2",
	}
	loginReq.ClientInfo = c.infoHandler.GetInfo()
	if err := c.WriteMessage(packet.MessageType_Login, loginReq); err != nil {
		log.Println("登录包发送失败:", err)
		return
	}
	log.Println("登录包发送成功")
	for {

	}
}

func (c *Client) ping() {
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
			if err := c.WriteMessage(packet.MessageType_Ping, &packet.Ping{}); err != nil {
				log.Println("心跳包发送失败:", err)
				return
			}
		}
	}

}

func (c *Client) WriteMessage(msgType packet.MessageType, msg proto.Message) error {
	body, err := anypb.New(msg)
	if err != nil {
		return err
	}

	data, err := proto.Marshal(&packet.Message{
		Body: body,
	})
	if err != nil {
		return err
	}
	return c.conn.WriteProtoMessage(msgType, data)
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
