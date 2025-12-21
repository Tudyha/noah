package session

import (
	"errors"
	"fmt"
	"io"
	"net"
	"noah/pkg/conn"
	"noah/pkg/logger"
	"noah/pkg/packet"
	"noah/pkg/utils"
	"sync/atomic"
	"time"

	smux "github.com/xtaci/smux/v2"
)

type Session struct {
	ID          uint64        // session id
	status      atomic.Int32  // session状态: 0: 初始化 1: 待认证 2: 认证成功 3: 关闭
	smuxSession *smux.Session // smux session

	conn *conn.Conn // 实际底层tcp
}

func newSession(netConn net.Conn) (*Session, error) {
	c := conn.NewConn(netConn)

	s := Session{
		ID:     uint64(utils.GenID()),
		status: atomic.Int32{},
		conn:   c,
	}

	// 设置session状态为待认证
	s.status.Store(1)

	// 设置认证超时，关闭无效连接
	authTime := 10 * time.Second
	logger.Info("session创建成功后，需在规定时间内完成认证，否则关闭连接", "sessionID", s.ID, "authTime", authTime)
	time.AfterFunc(authTime, func() {
		if s.status.Load() != 2 {
			logger.Error("session未完成认证，关闭连接", "sessionID", s.ID)
			s.Close()
		}
	})

	go s.checkAuth()
	return &s, nil
}

func (s *Session) checkAuth() error {
	var err error
	defer func() {
		if err != nil {
			logger.Error("session认证失败", "sessionID", s.ID, "err", err)
			s.Close()
		}
	}()

	p, err := s.conn.ReadOnce()
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
	if err = h.Handle(conn.NewConnContext(s.conn, p)); err != nil {
		return err
	}

	logger.Info("认证成功", "sessionID", s.ID)
	logger.Info("send login ack", "sessionID", s.ID)
	if _, err = s.conn.WriteProtoMessage(packet.MessageType_LoginAck, &packet.LoginAck{}); err != nil {
		return err
	}
	logger.Info("创建smux.session，并监听新连接", "sessionID", s.ID)

	smuxSession, err := smux.Server(s.conn.GetConn(), smux.DefaultConfig())
	if err != nil {
		logger.Error("创建smux session失败", "err", err)
		return err
	}
	s.smuxSession = smuxSession
	// 释放底层tcp，完全交由smux管理
	s.conn.Stop()
	s.conn = nil
	s.status.Store(2)
	return err
}

func (s *Session) accept() {
	for {
		c, err := s.smuxSession.AcceptStream()
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				return
			}
			continue
		}
		go s.handleConn(c)
	}
}

func (s *Session) handleConn(netConn net.Conn) {
	defer func() {
		netConn.Close()
	}()

	c := conn.NewConn(netConn)
	go c.Run()
	for {
		p, err := c.ReadMessage()
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				return
			}
			continue
		}
		ctx := conn.NewConnContext(c, p)
		h := messageHandlers[p.MessageType]
		if h == nil {
			logger.Error("消息处理函数未注册: %d", "messageType", p.MessageType)
			continue
		}
		if err = h.Handle(ctx); err != nil {
			logger.Error("处理消息失败: %s", "err", err)
		}
	}
}

func (s *Session) Close() error {
	if s.smuxSession != nil {
		s.smuxSession.Close()
	}
	if s.conn != nil {
		s.conn.Close()
	}
	s.status.Store(3)
	return nil
}
