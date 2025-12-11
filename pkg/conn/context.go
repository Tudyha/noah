package conn

import (
	"context"
	"noah/pkg/packet"
	"sync"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Context interface {
	context.Context
	GetConn() *Conn                           // 获取连接
	ShouldBindProto(body proto.Message) error // 解析消息体
	Release()                                 // 回收context
}

type MessageHandler interface {
	Handle(ctx Context) error
	MessageType() packet.MessageType
}

var pool = sync.Pool{
	New: func() any {
		return &connContext{}
	},
}

type connContext struct {
	context.Context
	conn      *Conn
	protoBody *anypb.Any
}

func NewConnContext(conn *Conn, p *packet.Packet) Context {
	c := pool.Get().(*connContext)
	c.conn = conn
	c.Context = context.Background()
	switch p.CodecType {
	case packet.CodecType_Protobuf:
		var msg packet.Message
		err := proto.Unmarshal(p.Data, &msg)
		if err == nil {
			c.protoBody = msg.Body
		}
	default:
	}
	return c
}

func (c *connContext) Release() {
	c.conn = nil
	c.protoBody = nil
	c.Context = nil
	pool.Put(c)
}

func (c *connContext) GetConn() *Conn {
	return c.conn
}

func (c *connContext) ShouldBindProto(body proto.Message) error {
	if c.protoBody == nil {
		return nil
	}
	return c.protoBody.UnmarshalTo(body)
}
