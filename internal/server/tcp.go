package server

import (
	"context"
	"errors"
	"net"
	"noah/internal/server/handler"
	"noah/pkg/app"
	"noah/pkg/config"
	"noah/pkg/conn"
	"noah/pkg/logger"
	"noah/pkg/packet"
	"sync"
	"time"
)

type tcpServer struct {
	listenAddr string       // 监听地址
	timeout    int          // 连接超时时间
	keepAlive  int          // 连接保活时间
	listener   net.Listener // 监听句柄，可用于关闭连接等

	sessions sync.Map // 连接会话，key -> connId, value -> *conn.Conn

	messageHandlers map[packet.MessageType]conn.MessageHandler
}

func NewTCPServer() app.Server {
	cfg := config.Get().Server.TCP
	s := &tcpServer{
		listenAddr: cfg.Addr,
		timeout:    cfg.Timeout,
		keepAlive:  cfg.KeepAlive,

		sessions: sync.Map{},

		messageHandlers: make(map[packet.MessageType]conn.MessageHandler),
	}

	registerHandler := handler.NewLoginHandler()
	pingHandler := handler.NewPingHandler()
	s.messageHandlers[registerHandler.MessageType()] = registerHandler
	s.messageHandlers[pingHandler.MessageType()] = pingHandler

	return s
}

// Start 启动TCP服务器
// 参数:
//
//	ctx: 上下文对象，用于控制监听操作的生命周期
//
// 返回值:
//
//	error: 启动过程中发生的错误，nil表示启动成功
func (t *tcpServer) Start(ctx context.Context) error {
	logger.Info("启动TCP服务", "listenAddr", t.listenAddr)

	// 配置监听器参数，设置KeepAlive时间
	listenerConfig := net.ListenConfig{
		KeepAlive: time.Duration(t.keepAlive) * time.Second,
	}

	// 创建TCP监听器
	listener, err := listenerConfig.Listen(ctx, "tcp", t.listenAddr)
	if err != nil {
		logger.Error("启动TCP服务失败", "err", err)
		return err
	}

	t.listener = listener

	// 启动监听协程处理连接
	go t.listen()

	return nil
}

// listen 监听TCP连接请求并处理新建立的连接
// 该方法会持续监听新的TCP连接，当有新连接建立时，启动一个新的goroutine来处理该连接
// 当监听器被关闭时，方法会正常退出
func (t *tcpServer) listen() {
	// 持续监听新的连接请求
	for {
		// 接受新的TCP连接
		conn, err := t.listener.Accept()
		if err != nil {
			// 检查是否因为监听器被关闭导致的错误
			if errors.Is(err, net.ErrClosed) {
				logger.Info("tcp server已停止, 结束连接监听")
				break
			}
			// 其他错误情况下继续监听
			continue
		}
		// 为每个新连接启动一个独立的goroutine进行处理
		go t.handleTCPConnection(conn)
	}
}

// handleTCPConnection 处理TCP连接的主循环函数
// 该函数负责接收和处理来自客户端的TCP连接，包括消息读取、解析和分发给相应的处理器
// 参数:
//   - netConn: 网络连接对象，用于与客户端进行数据通信
func (t *tcpServer) handleTCPConnection(netConn net.Conn) {
	c := conn.NewConn(netConn)
	defer c.Close()

	// 将新建立的连接添加到服务器的会话存储中，以便后续管理和查找
	t.sessions.Store(c.GetID(), c)

	// 主循环：持续读取消息并进行处理
	for {
		// 从连接中读取下一个消息
		p, err := c.ReadMessage()
		if err != nil {
			// 如果连接已关闭，则退出循环
			if errors.Is(err, net.ErrClosed) {
				break
			}
			// 其他错误情况下继续循环，尝试读取下一条消息
			continue
		}

		msgType := p.MessageType

		logger.Info("tcp server接收到消息", "msgType", msgType)

		// 检查连接状态，如果连接未激活且消息不是登录类型，则拒绝处理
		if c.GetState() != conn.ConnState_Active && msgType != packet.MessageType_Login {
			logger.Error("连接未激活，不允许处理消息", "connId", c.GetID())
			continue
		}

		// 创建上下文对象，用于在消息处理过程中传递连接相关信息
		ctx := conn.NewConnContext(c, p)

		// 根据消息类型查找对应的消息处理器
		h := t.messageHandlers[msgType]
		if h != nil {
			// 调用处理器处理消息
			if err := h.Handle(ctx); err != nil {
				logger.Error("msg handle err", "err", err)
			}
		}

		// 回收上下文对象
		ctx.Release()
	}
}

// Stop 停止TCP服务器
// 参数:
//
//	ctx - 上下文对象，用于控制停止过程
//
// 返回值:
//
//	error - 停止过程中可能产生的错误，当前实现始终返回nil
func (t *tcpServer) Stop(ctx context.Context) error {
	// 关闭监听器以停止接受新的连接
	if t.listener != nil {
		t.listener.Close()
	}
	return nil
}

func (t *tcpServer) String() string {
	return "TCP Server: " + t.listenAddr
}
