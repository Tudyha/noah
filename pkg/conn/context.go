package conn

import (
	"context"
	"noah/pkg/packet"
	"sync"

	"google.golang.org/protobuf/proto"
)

type Context interface {
	Close() error // 关闭连接
	Release()
}

type MessageHandler interface {
	Handle(ctx Context, msg proto.Message) error
	MessageType() packet.MessageType
	MessageBody() proto.Message
}

var pool = sync.Pool{
	New: func() any {
		return &connContext{}
	},
}

type connContext struct {
	context.Context
	conn *Conn
}

func NewConnContext(conn *Conn) Context {
	c := pool.Get().(*connContext)
	c.reset(conn)
	return c
}

func (c *connContext) Release() {
	c.conn = nil
	pool.Put(c)
}

func (c *connContext) reset(conn *Conn) {
	c.conn = conn
}

func (c *connContext) Close() error {
	return c.conn.Close()
}
