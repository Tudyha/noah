package session

import (
	"fmt"
	"net"
	"noah/internal/handler"
	"noah/pkg/conn"
	"noah/pkg/logger"
	"noah/pkg/packet"
	"sync"

	"google.golang.org/protobuf/proto"
)

type SessionManager interface {
	NewSession(netConn net.Conn) (*Session, error)
	handleMessage(ctx conn.Context, msgType packet.MessageType) error
	SendProtoMessage(connID uint64, msgType packet.MessageType, msg proto.Message) error
}

var (
	sessionManagerInstance SessionManager
)

type sessionManager struct {
	sessions        sync.Map                                   // 会话, key -> connID, value -> *Session
	messageHandlers map[packet.MessageType]conn.MessageHandler // 消息处理器
}

func Init() error {
	sm := &sessionManager{
		sessions:        sync.Map{},
		messageHandlers: make(map[packet.MessageType]conn.MessageHandler),
	}
	registerHandler := handler.NewLoginHandler()
	pingHandler := handler.NewPingHandler()
	sm.messageHandlers[registerHandler.MessageType()] = registerHandler
	sm.messageHandlers[pingHandler.MessageType()] = pingHandler

	sessionManagerInstance = sm
	return nil
}

func GetSessionManager() SessionManager {
	return sessionManagerInstance
}

func (m *sessionManager) NewSession(netConn net.Conn) (*Session, error) {
	logger.Info("create client session")
	c := conn.NewConn(netConn)

	s := Session{
		conn: c,
	}

	m.sessions.Store(c.GetID(), &s)

	go s.readMessage()
	return &s, nil
}

func (m *sessionManager) handleMessage(ctx conn.Context, msgType packet.MessageType) error {
	h, ok := m.messageHandlers[msgType]
	if !ok {
		return fmt.Errorf("msg handler not found, msgType: %d", msgType)
	}
	return h.Handle(ctx)
}

func (m *sessionManager) SendProtoMessage(connID uint64, msgType packet.MessageType, msg proto.Message) error {
	s, ok := m.sessions.Load(connID)
	if !ok {
		return fmt.Errorf("session not found, connID: %d", connID)
	}
	return s.(*Session).SendProtoMessage(msgType, msg)
}
