package session

import (
	"errors"
	"io"
	"net"
	"noah/pkg/conn"
	"noah/pkg/logger"
	"noah/pkg/packet"

	"google.golang.org/protobuf/proto"
)

type Session struct {
	conn *conn.Conn
}

func (s *Session) readMessage() {
	for {
		// 从连接中读取下一个消息
		p, err := s.conn.ReadMessage()
		if err != nil {
			// 如果连接已关闭，则退出循环
			if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
				break
			}
			// 其他错误情况下继续循环，尝试读取下一条消息
			continue
		}

		msgType := p.MessageType

		logger.Info("tcp server接收到消息", "msgType", msgType)

		// 检查连接状态，如果连接未激活且消息不是登录类型，则拒绝处理
		if s.conn.GetState() != conn.ConnState_Active && msgType != packet.MessageType_Login {
			logger.Error("连接未激活，不允许处理消息", "connId", s.conn.GetID())
			continue
		}

		// 创建上下文对象，用于在消息处理过程中传递连接相关信息
		ctx := conn.NewConnContext(s.conn, p)

		if err = sessionManagerInstance.handleMessage(ctx, msgType); err != nil {
			logger.Error("msg handle err", "err", err)
		}

		// 回收上下文对象
		ctx.Release()
	}
}

func (s *Session) SendProtoMessage(msgType packet.MessageType, msg proto.Message) error {
	return s.conn.WriteProtoMessage(msgType, msg)
}
