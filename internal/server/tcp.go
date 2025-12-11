package server

import (
	"context"
	"errors"
	"net"
	"noah/internal/server/handler"
	"noah/internal/service"
	"noah/pkg/app"
	"noah/pkg/config"
	"noah/pkg/conn"
	"noah/pkg/logger"
	"noah/pkg/packet"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type tcpServer struct {
	listenAddr    string                // 监听地址
	timeout       int                   // 连接超时时间
	keepAlive     int                   // 连接保活时间
	listener      net.Listener          // 监听句柄，可用于关闭连接等
	clientService service.ClientService // 客户端服务，处理客户端逻辑

	sessions sync.Map // smux 会话

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
	s.messageHandlers[registerHandler.MessageType()] = registerHandler

	return s
}

func (t *tcpServer) Start(ctx context.Context) error {
	logger.Info("启动TCP服务", "listenAddr", t.listenAddr)
	listenerConfig := net.ListenConfig{
		KeepAlive: time.Duration(t.keepAlive) * time.Second,
	}
	listener, err := listenerConfig.Listen(ctx, "tcp", t.listenAddr)
	if err != nil {
		logger.Error("启动TCP服务失败", "err", err)
		return err
	}
	t.listener = listener
	go t.listen()
	return nil
}

func (t *tcpServer) listen() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				logger.Info("tcp server已停止, 结束连接监听")
				break
			}
			continue
		}
		go t.handleTCPConnection(conn)
	}
}

func (t *tcpServer) handleTCPConnection(netConn net.Conn) {
	c := conn.NewConn(netConn)
	defer c.Close()

	ctx := conn.NewConnContext(c)

	for {
		msgType, msg, err := c.ReadMessage()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			continue
		}
		logger.Info("tcp server接收到消息", "msgType", msgType, "msg", msg)
		if c.GetState() != conn.ConnState_Active && msgType != packet.MessageType_Login {
			logger.Error("连接未激活，不允许处理消息", "connId", c.GetID())
			continue
		}

		h := t.messageHandlers[msgType]
		if h != nil {
			body := h.MessageBody()
			if err := anypb.UnmarshalTo(msg, body, proto.UnmarshalOptions{}); err != nil {
				logger.Error("anypb.UnmarshalTo err", "err", err)
				continue
			}
			if err := h.Handle(ctx, body); err != nil {
				logger.Error("msg handle err", "err", err)
			}
		}
	}
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
