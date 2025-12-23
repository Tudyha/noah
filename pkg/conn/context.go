package conn

import (
	"context"
	"io"
	"noah/pkg/packet"
	"sync"

	"google.golang.org/protobuf/proto"
)

type Context interface {
	context.Context
	GetConn() *Conn                           // 获取连接
	ShouldBindProto(body proto.Message) error // 解析消息体
	Release()                                 // 回收context
	WithValue(key any, value any)
	IsHijacked() bool // 是否被劫持
	Hijack() (io.ReadWriteCloser, error)
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
	conn       *Conn
	request    *packet.Packet
	isHijacked bool
}

func NewConnContext(conn *Conn, p *packet.Packet) Context {
	c := pool.Get().(*connContext)
	c.conn = conn
	c.Context = context.Background()
	c.request = p
	return c
}

func (c *connContext) Release() {
	c.conn = nil
	c.request = nil
	c.Context = nil
	pool.Put(c)
}

func (c *connContext) GetConn() *Conn {
	return c.conn
}

func (c *connContext) ShouldBindProto(body proto.Message) error {
	if c.request == nil {
		return nil
	}
	return c.request.Unmarshal(body)
}

func (c *connContext) WithValue(key any, value any) {
	c.Context = context.WithValue(c.Context, key, value)
}

func (c *connContext) Hijack() (io.ReadWriteCloser, error) {
	c.isHijacked = true
	rwc := c.conn
	return rwc, nil
}

func (c *connContext) IsHijacked() bool {
	return c.isHijacked
}
