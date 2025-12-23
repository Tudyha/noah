package session

import (
	"io"
	"net"
	"noah/internal/handler"
	"noah/pkg/conn"
	"noah/pkg/errcode"
	"noah/pkg/logger"
	"noah/pkg/packet"
	"sync"

	"google.golang.org/protobuf/proto"
)

type SessionManager interface {
	NewSession(netConn net.Conn) (*Session, error)
	SendProtoMessage(sessionID uint64, msgType packet.MessageType, msg proto.Message) error
	OpenTunnel(sessionID uint64, tunnelType packet.OpenTunnel_TuunnelType) (io.ReadWriteCloser, error)
}

var (
	sessionManagerInstance SessionManager
	messageHandlers        map[packet.MessageType]conn.MessageHandler // 消息处理器

	localSessionManager *sessionManager
)

type sessionManager struct {
	sessions sync.Map // 会话, key -> sessionID, value -> *Session
}

func Init() error {
	sm := &sessionManager{
		sessions: sync.Map{},
	}

	// 注册消息处理器
	messageHandlers = make(map[packet.MessageType]conn.MessageHandler)
	registerHandler := handler.NewLoginHandler()
	pingHandler := handler.NewPingHandler()
	messageHandlers[registerHandler.MessageType()] = registerHandler
	messageHandlers[pingHandler.MessageType()] = pingHandler

	sessionManagerInstance = sm
	localSessionManager = sm
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

	m.sessions.Store(s.ID, s)

	return s, nil
}

func (m *sessionManager) SendProtoMessage(sessionID uint64, msgType packet.MessageType, msg proto.Message) error {
	v, ok := m.sessions.Load(sessionID)
	if !ok {
		return errcode.ErrClientDisconnect
	}
	s, ok := v.(*Session)
	if !ok {
		return errcode.ErrClientDisconnect
	}
	if s.status.Load() != 2 {
		return errcode.ErrClientDisconnect
	}
	return s.WriteProtoMessage(msgType, msg)
}

func (m *sessionManager) OpenTunnel(sessionID uint64, tunnelType packet.OpenTunnel_TuunnelType) (io.ReadWriteCloser, error) {
	v, ok := m.sessions.Load(sessionID)
	if !ok {
		return nil, errcode.ErrClientDisconnect
	}
	s, ok := v.(*Session)
	if !ok {
		return nil, errcode.ErrClientDisconnect
	}
	if s.status.Load() != 2 {
		return nil, errcode.ErrClientDisconnect
	}
	return s.OpenTunnel(tunnelType)
}
