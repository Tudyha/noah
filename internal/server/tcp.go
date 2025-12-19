package server

import (
	"context"
	"errors"
	"net"
	"noah/internal/session"
	"noah/pkg/app"
	"noah/pkg/config"
	"noah/pkg/logger"
	"time"
)

type tcpServer struct {
	listenAddr string       // 监听地址
	timeout    int          // 连接超时时间
	keepAlive  int          // 连接保活时间
	listener   net.Listener // 监听句柄，可用于关闭连接等

	sessionManager session.SessionManager
}

func NewTCPServer() app.Server {
	cfg := config.Get().Server.TCP
	s := &tcpServer{
		listenAddr:     cfg.Addr,
		timeout:        cfg.Timeout,
		keepAlive:      cfg.KeepAlive,
		sessionManager: session.GetSessionManager(),
	}

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
		logger.Info("new tcp conn", "RemoteAddr", conn.RemoteAddr().String())
		_, err = t.sessionManager.NewSession(conn)
		if err != nil {
			logger.Error("tcp server创建session失败", "err", err)
			conn.Close()
		}
	}

	return nil
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
