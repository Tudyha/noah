package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"noah/client/app/handler"
	pkgApp "noah/pkg/app"
	"noah/pkg/config"
	"noah/pkg/conn"
	"noah/pkg/packet"
	"noah/pkg/utils"
	"sync"
	"sync/atomic"
	"time"

	"github.com/xtaci/smux/v2"
)

type Client struct {
	cfg *config.ClientConfig

	connected atomic.Bool

	// Use RWMutex to protect session-related resources
	mu         sync.RWMutex
	session    *smux.Session
	pingStream *smux.Stream

	// Lifecycle control
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	infoHandler     *handler.InfoHandler
	messageHandlers map[packet.MessageType]conn.MessageHandler
}

func NewClient(cfg *config.ClientConfig) pkgApp.Server {
	// Set defaults
	if cfg.ReconnectInterval <= 0 {
		cfg.ReconnectInterval = 5
	}
	if cfg.DailTimeout <= 0 {
		cfg.DailTimeout = 10
	}
	if cfg.HeartbeatInterval <= 0 {
		cfg.HeartbeatInterval = 30
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := &Client{
		cfg:             cfg,
		ctx:             ctx,
		cancel:          cancel,
		infoHandler:     &handler.InfoHandler{},
		messageHandlers: make(map[packet.MessageType]conn.MessageHandler),
	}

	// Registry handlers
	logoutHandler := handler.NewLogoutHandler()
	tunnelHandler := handler.NewTunnelHandler()
	c.messageHandlers[logoutHandler.MessageType()] = logoutHandler
	c.messageHandlers[tunnelHandler.MessageType()] = tunnelHandler

	return c
}

func (c *Client) Start(ctx context.Context) error {
	log.Println("Starting client...")

	// Initial connection attempt
	go c.reconnectLoop()

	// Wait for global context or local cancel
	select {
	case <-ctx.Done():
		return c.Stop(ctx)
	case <-c.ctx.Done():
		return nil
	}
}

func (c *Client) reconnectLoop() {
	backoff := time.Duration(c.cfg.ReconnectInterval) * time.Second
	maxBackoff := 60 * time.Second

	for {
		if !c.connected.Load() {
			if err := c.connect(); err != nil {
				log.Printf("Connection failed: %v, retrying in %v...", err, backoff)

				// Exponential backoff
				select {
				case <-c.ctx.Done():
					return
				case <-time.After(backoff):
					backoff *= 2
					if backoff > maxBackoff {
						backoff = maxBackoff
					}
					continue
				}
			}
			// Reset backoff on success
			backoff = time.Duration(c.cfg.ReconnectInterval) * time.Second
		}

		select {
		case <-c.ctx.Done():
			return
		case <-time.After(5 * time.Second): // Health check interval
		}
	}
}

func (c *Client) connect() error {
	dialer := net.Dialer{
		Timeout: time.Duration(c.cfg.DailTimeout) * time.Second,
	}

	netConn, err := dialer.DialContext(c.ctx, "tcp", c.cfg.ServerAddr)
	if err != nil {
		return err
	}

	// Start auth and session promotion
	go c.handleConn(netConn)
	return nil
}

func (c *Client) handleConn(netConn net.Conn) {
	co := conn.NewConn(netConn)

	// Local cleanup logic
	var sess *smux.Session
	var ping *smux.Stream

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in handleConn: %v", r)
		}
		if !c.connected.Load() {
			if co != nil {
				co.Close()
			}
			if sess != nil {
				sess.Close()
			}
			netConn.Close()
		}
	}()

	loginReq := &packet.Login{
		AppId:    c.cfg.AppId,
		Sign:     utils.Sign(c.cfg.AppId, c.cfg.AppSecret),
		DeviceId: utils.GetMacAddress(),
	}
	loginReq.ClientInfo = c.infoHandler.GetInfo()

	co.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err := co.WriteProtoMessage(packet.MessageType_Login, loginReq); err != nil {
		return
	}

	co.SetReadDeadline(time.Now().Add(10 * time.Second))
	p, err := co.ReadMessage()
	if err != nil || p.MessageType != packet.MessageType_LoginAck {
		log.Printf("Auth failed: %v", err)
		return
	}
	co.SetDeadline(time.Time{}) // Clear deadlines

	log.Println("Auth passed, upgrading to smux")
	co.Release()

	// 2. SMUX Promotion
	sess, err = smux.Client(netConn, smux.DefaultConfig())
	if err != nil {
		return
	}

	// 3. Heartbeat setup
	ping, err = sess.OpenStream()
	if err != nil {
		sess.Close()
		return
	}

	// 4. Critical Section: Update Client State
	c.mu.Lock()
	c.session = sess
	c.pingStream = ping
	c.mu.Unlock()

	c.connected.Store(true)

	// 5. Spin up workers
	c.wg.Add(2)
	go c.heartbeatWorker(ping)
	go c.acceptStreamLoop(sess)
}

func (c *Client) heartbeatWorker(s *smux.Stream) {
	defer c.wg.Done()
	defer c.markDisconnected()

	co := conn.NewConn(s)
	ticker := time.NewTicker(time.Duration(c.cfg.HeartbeatInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			data := c.infoHandler.GetSystemStat()
			if err := co.WriteProtoMessage(packet.MessageType_Ping, data); err != nil {
				log.Printf("Heartbeat failed: %v", err)
				return
			}
		}
	}
}

func (c *Client) acceptStreamLoop(sess *smux.Session) {
	defer c.wg.Done()
	defer c.markDisconnected()

	for {
		stream, err := sess.AcceptStream()
		if err != nil {
			return
		}
		go c.handleStream(stream)
	}
}

func (c *Client) handleStream(netConn net.Conn) {
	co := conn.NewConn(netConn)

	for {
		p, err := co.ReadMessage()
		if err != nil {
			co.Close()
			return
		}

		ctx := conn.NewConnContext(co, p)
		h := c.messageHandlers[p.MessageType]
		if h == nil {
			log.Printf("Handler missing for type: %v", p.MessageType)
			ctx.Release()
			continue
		}

		if err := h.Handle(ctx); err != nil {
			log.Printf("Handle error: %v", err)
		}

		hijacked := ctx.IsHijacked()
		ctx.Release()
		if hijacked {
			return
		}
	}
}

func (c *Client) markDisconnected() {
	if c.connected.Swap(false) {
		log.Println("Client disconnected, cleaning up...")
		c.mu.Lock()
		if c.session != nil {
			c.session.Close()
		}
		c.mu.Unlock()
	}
}

func (c *Client) Stop(ctx context.Context) error {
	log.Println("Stopping client...")
	c.cancel() // Signal all workers to stop

	c.mu.Lock()
	if c.session != nil {
		c.session.Close()
	}
	c.mu.Unlock()

	// Wait for goroutines to finish (graceful)
	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("Client stopped gracefully")
	case <-ctx.Done():
		log.Println("Stop timed out, forcing exit")
	}
	return nil
}

func (c *Client) String() string {
	return fmt.Sprintf("NoahClient[%d]", c.cfg.AppId)
}
