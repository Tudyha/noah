package session

import (
	"errors"
	"fmt"
	"io"
	"net"
	"noah/pkg/conn"
	"noah/pkg/constant"
	"noah/pkg/logger"
	"noah/pkg/packet"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	smux "github.com/xtaci/smux/v2"
)

// Session Status Constants
const (
	StatusInit    int32 = 0
	StatusPending int32 = 1 // Awaiting authentication
	StatusReady   int32 = 2 // Authenticated and smux ready
	StatusClosed  int32 = 3
)

var (
	ErrSessionNotReady = errors.New("session not authenticated or ready")
)

// 客户端会话
type Session struct {
	ID          string        // session id
	status      atomic.Int32  // session状态: 0: 初始化 1: 待认证 2: 认证成功 3: 关闭
	smuxSession *smux.Session // smux多路复用器

	closeOnce sync.Once
}

// 创建新会话
func newSession(netConn net.Conn) (*Session, error) {
	s := &Session{
		ID: uuid.NewString(),
	}

	// 设置session状态为待认证
	s.status.Store(StatusPending)

	// 检测认证
	go s.checkAuth(netConn)
	return s, nil
}

// 鉴权
func (s *Session) checkAuth(netConn net.Conn) {
	var err error
	c := conn.NewConn(netConn)

	defer func() {
		if err != nil || s.status.Load() != StatusReady {
			logger.Error("session认证失败", "sessionID", s.ID, "err", err)
			c.Close()
			s.Close()
		}
	}()

	// 读取鉴权包 TODO 增加读取超时控制
	c.SetReadDeadline(time.Now().Add(10 * time.Second))
	p, err := c.ReadMessage()
	if err != nil {
		return
	}
	c.SetDeadline(time.Time{})
	if p.MessageType != packet.MessageType_Login {
		err = fmt.Errorf("认证失败，消息类型错误: %d", p.MessageType)
		return
	}
	h := messageHandlers[p.MessageType]
	if h == nil {
		err = fmt.Errorf("认证失败，消息处理函数未注册: %d", p.MessageType)
		return
	}
	// 校验鉴权信息
	msgContext := conn.NewConnContext(c, p)
	msgContext.WithValue(constant.SESSION_ID_KEY, s.ID)
	err = h.Handle(msgContext)
	msgContext.Release()
	if err != nil {
		return
	}

	logger.Info("认证成功", "sessionID", s.ID)
	if err = c.WriteProtoMessage(packet.MessageType_LoginAck, &packet.LoginAck{}); err != nil {
		return
	}

	// 移交net.Conn给smux前先释放当前持有
	c.Release()
	c = nil

	logger.Info("创建smux session", "sessionID", s.ID)
	smuxConfig := smux.DefaultConfig()
	smuxSession, err := smux.Server(netConn, smuxConfig)
	if err != nil {
		return
	}
	s.smuxSession = smuxSession

	// 设置session状态为已认证
	s.status.Store(StatusReady)

	go s.accept()
}

// 接受新子流
func (s *Session) accept() {
	defer s.Close()
	for {
		c, err := s.smuxSession.AcceptStream()
		if err != nil {
			logger.Error("接受子流失败", "sessionID", s.ID, "err", err)
			return
		}
		go s.handleStream(c)
	}
}

// 处理子流
func (s *Session) handleStream(netConn net.Conn) {
	c := conn.NewConn(netConn)

	for {
		p, err := c.ReadMessage()
		if err != nil {
			logger.Error("读取子流消息失败", "err", err)
			c.Close()
			return
		}
		// 创建上下文，用于传递数据
		ctx := conn.NewConnContext(c, p)
		ctx.WithValue(constant.SESSION_ID_KEY, s.ID)
		h := messageHandlers[p.MessageType]
		if h == nil {
			logger.Warn("消息处理函数未注册", "messageType", p.MessageType)
			ctx.Release()
			continue
		}
		if err = h.Handle(ctx); err != nil {
			logger.Error("处理消息失败", "err", err)
		}

		hijacked := ctx.IsHijacked()
		ctx.Release()

		if hijacked {
			return
		}
	}
}

// 关闭会话
func (s *Session) Close() error {
	s.closeOnce.Do(func() {
		logger.Info("close session", "sessionID", s.ID)
		s.status.Store(StatusClosed)

		if s.smuxSession != nil {
			s.smuxSession.Close()
		}

		localSessionManager.sessions.Delete(s.ID)
	})
	return nil
}

func (s *Session) isReady() (*smux.Session, error) {
	if s.status.Load() != StatusReady {
		return nil, ErrSessionNotReady
	}
	if s.smuxSession == nil {
		return nil, ErrSessionNotReady
	}
	return s.smuxSession, nil
}

// 发送proto消息
func (s *Session) SendCommand(cmd packet.Command_Cmd) error {
	// fixme: 临时方案，考虑增加专门用于发送消息的stream
	smuxSession, err := s.isReady()
	if err != nil {
		return err
	}
	netConn, err := smuxSession.OpenStream()
	if err != nil {
		return err
	}
	logger.Info("create steam success", "sessionID", s.ID)
	defer netConn.Close()
	c := conn.NewConn(netConn)
	defer c.Close()

	return c.WriteProtoMessage(packet.MessageType_Command, &packet.Command{
		Cmd: cmd,
	})
}

// 打开tunnel
func (s *Session) OpenTunnel(tunnelType packet.OpenTunnel_TuunnelType, addr string) (io.ReadWriteCloser, error) {
	var err error

	smuxSession, err := s.isReady()
	if err != nil {
		return nil, err
	}
	netConn, err := smuxSession.OpenStream()
	if err != nil {
		return nil, err
	}
	c := conn.NewConn(netConn)
	defer func() {
		if err != nil {
			netConn.Close()
			c.Close()
		}
	}()

	logger.Info("open tunnel", "tunnel_type", tunnelType, "addr", addr)
	if err = c.WriteProtoMessage(packet.MessageType_Tunnel_Open, &packet.OpenTunnel{
		TunnelType: tunnelType,
		Addr:       addr,
	}); err != nil {
		return nil, err
	}

	p, err := c.ReadMessage()
	if err != nil {
		return nil, err
	}
	if p.MessageType != packet.MessageType_Tunnel_Open_Ack {
		err = fmt.Errorf("tunnel open ack error: %d", p.MessageType)
		return nil, err
	}
	var ack packet.OpenTunnelAck
	if err = p.Unmarshal(&ack); err != nil {
		return nil, err
	}
	if ack.Code != 0 {
		err = fmt.Errorf("tunnel open ack error: %s", ack.Msg)
		return nil, err
	}

	return c, nil
}
