package conn

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"noah/pkg/packet"
	"noah/pkg/utils"
	"sync"
	"time"
)

type ConnState uint8

const (
	ConnState_Init ConnState = iota
	ConnState_Active
	ConnState_Closed
)

type Conn struct {
	connID  uint64   // 连接ID
	netConn net.Conn // 底层连接

	codec packet.Codec // 编解码器

	state      ConnState // 连接状态
	stateMutex sync.RWMutex

	readMutex  sync.Mutex    // 读锁
	writeMutex sync.Mutex    // 写锁
	muxBuf     *bytes.Buffer // 多路复用器数据缓冲区
	cond       *sync.Cond    // 用于在muxBuf为空时阻塞Read()，并在写入数据时唤醒

	messageChan chan *packet.Packet // 消息通道

	closeOnce sync.Once
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewConn(conn net.Conn) *Conn {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Conn{
		connID:      uint64(utils.GenID()),
		codec:       packet.NewCodec(),
		netConn:     conn,
		muxBuf:      new(bytes.Buffer),
		messageChan: make(chan *packet.Packet, 1024),
		ctx:         ctx,
		cancel:      cancel,
	}
	c.cond = sync.NewCond(&c.readMutex) // 初始化 Cond

	go c.read()

	return c
}

func (c *Conn) read() {
	defer c.Close()

	for c.ctx.Err() == nil {
		p, err := c.codec.Decode(c.netConn)
		if err != nil {
			fmt.Println(err)
			return
		}

		c.readMutex.Lock()
		switch p.MessageType {
		case packet.MessageType_Stream_Data:
			c.muxBuf.Write(p.Data)
			c.cond.Broadcast()
		default:
			c.messageChan <- p
		}
		c.readMutex.Unlock()
	}
}

func (c *Conn) Read(b []byte) (n int, err error) {
	c.readMutex.Lock()
	defer c.readMutex.Unlock()
	for {
		if c.muxBuf.Len() > 0 {
			return c.muxBuf.Read(b)
		}

		if c.ctx.Err() != nil {
			return 0, io.EOF
		}

		c.cond.Wait()
	}
}

func (c *Conn) Write(b []byte) (n int, err error) {
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()
	p := &packet.Packet{
		MessageType: packet.MessageType_Stream_Data,
		Data:        b,
	}

	return c.codec.Encode(c.netConn, p)
}

func (c *Conn) Close() error {
	c.closeOnce.Do(func() {
		c.cancel()

		c.cond.Broadcast()

		close(c.messageChan)

		if c.netConn != nil {
			c.netConn.Close()
		}
	})
	return nil
}

func (c *Conn) GetID() uint64 {
	return c.connID
}

func (c *Conn) ReadMessage() (*packet.Packet, error) {
	select {
	case <-c.ctx.Done():
		return nil, io.EOF
	case p, ok := <-c.messageChan:
		if !ok {
			return nil, io.EOF
		}
		return p, nil
	}
}

func (c *Conn) WriteProtoMessage(msgType packet.MessageType, data []byte) error {
	p := &packet.Packet{
		MessageType: msgType,
		CodecType:   packet.CodecType_Protobuf,
		Data:        data,
	}
	_, err := c.codec.Encode(c.netConn, p)
	return err
}

func (c *Conn) GetState() ConnState {
	c.stateMutex.RLock()
	defer c.stateMutex.RUnlock()
	return c.state
}

func (c *Conn) LocalAddr() net.Addr {
	if c.netConn != nil {
		return c.netConn.LocalAddr()
	}
	return nil
}

func (c *Conn) RemoteAddr() net.Addr {
	if c.netConn != nil {
		return c.netConn.RemoteAddr()
	}
	return nil
}

func (c *Conn) SetDeadline(t time.Time) error {
	if c.netConn != nil {
		return c.netConn.SetDeadline(t)
	}
	return nil
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	if c.netConn != nil {
		return c.netConn.SetReadDeadline(t)
	}
	return nil
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	if c.netConn != nil {
		return c.netConn.SetWriteDeadline(t)
	}
	return nil
}

func (c *Conn) SetState(state ConnState) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()
	c.state = state
}
