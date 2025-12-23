package session

import (
	"fmt"
	"io"
	"net"
	"noah/pkg/conn"
	"noah/pkg/constant"
	"noah/pkg/logger"
	"noah/pkg/packet"
	"noah/pkg/utils"
	"sync/atomic"

	smux "github.com/xtaci/smux/v2"
	"google.golang.org/protobuf/proto"
)

// 客户端会话
type Session struct {
	ID          uint64        // session id
	status      atomic.Int32  // session状态: 0: 初始化 1: 待认证 2: 认证成功 3: 关闭
	smuxSession *smux.Session // smux session smux多路复用器
}

// 创建新会话
func newSession(netConn net.Conn) (*Session, error) {
	s := Session{
		ID:     uint64(utils.GenID()),
		status: atomic.Int32{},
	}

	// 设置session状态为待认证
	s.status.Store(1)

	// 检测认证
	go s.checkAuth(netConn)
	return &s, nil
}

// 鉴权
func (s *Session) checkAuth(netConn net.Conn) error {
	c := conn.NewConn(netConn)
	var err error
	defer func() {
		if err != nil || s.status.Load() != 2 {
			logger.Error("session认证失败", "sessionID", s.ID, "err", err)
			netConn.Close()
			s.Close()
		}
	}()

	// 读取鉴权包 TODO 增加超时
	p, err := c.ReadMessage()
	if err != nil {
		return err
	}
	if p.MessageType != packet.MessageType_Login {
		err = fmt.Errorf("认证失败，消息类型错误: %d", p.MessageType)
		return err
	}
	h := messageHandlers[p.MessageType]
	if h == nil {
		err = fmt.Errorf("认证失败，消息处理函数未注册: %d", p.MessageType)
		return err
	}
	// 校验鉴权信息
	msgContext := conn.NewConnContext(c, p)
	msgContext.WithValue(constant.SESSION_ID_KEY, s.ID)
	if err = h.Handle(msgContext); err != nil {
		msgContext.Release()
		return err
	}
	msgContext.Release()

	logger.Info("认证成功", "sessionID", s.ID)
	if err = c.WriteProtoMessage(packet.MessageType_LoginAck, &packet.LoginAck{}); err != nil {
		return err
	}

	// 移交net.Conn给smux前先释放当前持有
	c.Release()
	c = nil

	logger.Info("创建smux.session，并监听新连接", "sessionID", s.ID)
	smuxSession, err := smux.Server(netConn, smux.DefaultConfig())
	if err != nil {
		logger.Error("创建smux session失败", "err", err)
		return err
	}
	s.smuxSession = smuxSession

	s.status.Store(2)

	go s.accept()
	return err
}

// 接受新子流
func (s *Session) accept() {
	defer s.Close()
	for {
		c, err := s.smuxSession.AcceptStream()
		if err != nil {
			return
		}
		go s.handleConn(c)
	}
}

// 处理子流
func (s *Session) handleConn(netConn net.Conn) {
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
			ctx.Release()
			logger.Error("消息处理函数未注册", "messageType", p.MessageType)
			continue
		}
		if err = h.Handle(ctx); err != nil {
			logger.Error("处理消息失败", "err", err)
		}

		if ctx.IsHijacked() {
			logger.Info("底层连接已被劫持，不再处理消息", "sessionID", s.ID)
			ctx.Release()
			return
		}
		ctx.Release()
	}
}

// 关闭会话
func (s *Session) Close() error {
	logger.Info("close session", "sessionID", s.ID)
	if s.smuxSession != nil {
		s.smuxSession.Close()
	}
	s.status.Store(3)

	localSessionManager.sessions.Delete(s.ID)
	return nil
}

// 发送proto消息
func (s *Session) WriteProtoMessage(msgType packet.MessageType, msg proto.Message) error {
	// fixme: 临时方案，考虑增加专门用于发送消息的stream
	netConn, err := s.smuxSession.OpenStream()
	if err != nil {
		return err
	}
	logger.Info("create steam success", "sessionID", s.ID)
	defer netConn.Close()
	c := conn.NewConn(netConn)
	defer c.Close()

	return c.WriteProtoMessage(msgType, msg)
}

func (s *Session) OpenTunnel(tunnelType packet.OpenTunnel_TuunnelType) (io.ReadWriteCloser, error) {
	var err error

	netConn, err := s.smuxSession.OpenStream()
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

	if err = c.WriteProtoMessage(packet.MessageType_Tunnel_Open, &packet.OpenTunnel{
		TunnelType: tunnelType,
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
