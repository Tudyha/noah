package server

import (
	"context"
	"errors"
	"net"
	"noah/internal/session"
	"noah/pkg/app"
	"noah/pkg/config"
	"noah/pkg/logger"
)

type tcpServer struct {
	listenAddr string       // 监听地址
	listener   net.Listener // 监听句柄，可用于关闭连接等

	sessionManager session.SessionManager // 会话管理器
}

func NewTCPServer() app.Server {
	cfg := config.Get().Server.TCP
	s := &tcpServer{
		listenAddr:     cfg.Addr,
		sessionManager: session.GetSessionManager(),
	}

	return s
}

func (t *tcpServer) Start(ctx context.Context) error {
	logger.Info("tcp server start", "addr", t.listenAddr)
	listener, err := net.Listen("tcp", t.listenAddr)
	if err != nil {
		logger.Error("启动TCP服务失败", "err", err)
		return err
	}

	t.listener = listener

	for {
		conn, err := t.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				logger.Info("tcp server已停止, 结束连接监听")
				break
			}
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

func (t *tcpServer) Stop(ctx context.Context) error {
	if t.listener != nil {
		t.listener.Close()
	}
	return nil
}

func (t *tcpServer) String() string {
	return "TCP Server: " + t.listenAddr
}
