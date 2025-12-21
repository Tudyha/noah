package session

import (
	"net"
	"noah/internal/handler"
	"noah/pkg/conn"
	"noah/pkg/logger"
	"noah/pkg/packet"
	"sync"
)

type SessionManager interface {
	NewSession(netConn net.Conn) (*Session, error)
}

var (
	sessionManagerInstance SessionManager
	messageHandlers        map[packet.MessageType]conn.MessageHandler // 消息处理器
)

type sessionManager struct {
	sessions sync.Map // 会话, key -> connID, value -> *Session
}

func Init() error {
	sm := &sessionManager{
		sessions: sync.Map{},
	}

	messageHandlers = make(map[packet.MessageType]conn.MessageHandler)
	registerHandler := handler.NewLoginHandler()
	pingHandler := handler.NewPingHandler()
	messageHandlers[registerHandler.MessageType()] = registerHandler
	messageHandlers[pingHandler.MessageType()] = pingHandler

	sessionManagerInstance = sm
	return nil
}

func GetSessionManager() SessionManager {
	return sessionManagerInstance
}

func (m *sessionManager) NewSession(netConn net.Conn) (*Session, error) {
	logger.Info("create conn session")

	s, err := newSession(netConn)
	if err != nil {
		return nil, err
	}

	m.sessions.Store(s.ID, &s)

	return s, nil
}
