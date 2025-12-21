package conn

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"noah/pkg/packet"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Conn struct {
	netConn     net.Conn     // 底层连接
	codec       packet.Codec // 编解码器
	closeOnce   sync.Once
	closeSignal chan struct{}

	messageChan chan *packet.Packet // 控制消息
	dataBuf     *bytes.Buffer       // 数据
	cond        *sync.Cond
	readMux     *sync.Mutex
	writeMux    *sync.Mutex
}

func NewConn(conn net.Conn) *Conn {
	c := &Conn{
		codec:       packet.NewCodec(),
		netConn:     conn,
		closeSignal: make(chan struct{}),
		messageChan: make(chan *packet.Packet, 1024),
		dataBuf:     bytes.NewBuffer(nil),
		readMux:     &sync.Mutex{},
		writeMux:    &sync.Mutex{},
	}
	c.cond = sync.NewCond(c.readMux)

	return c
}

func (c *Conn) Run() {
	defer func() {
		c.Close()
	}()
	for {
		select {
		case <-c.closeSignal:
			fmt.Println("stop read")
			return
		default:
			p, err := c.codec.Decode(c.netConn)
			if err != nil {
				if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
					return
				}
				continue
			}
			c.readMux.Lock()
			switch p.MessageType {
			case packet.MessageType_Stream_Data:
				c.dataBuf.Write(p.Data)
				c.cond.Broadcast()
			default:
				c.messageChan <- p
			}
			c.readMux.Unlock()
		}
	}
}

func (c *Conn) ReadOnce() (*packet.Packet, error) {
	return c.codec.Decode(c.netConn)
}

func (c *Conn) ReadMessage() (*packet.Packet, error) {
	p, ok := <-c.messageChan
	if !ok {
		return nil, io.EOF
	}
	return p, nil
}

func (c *Conn) WriteProtoMessage(msgType packet.MessageType, msg proto.Message) (int, error) {
	c.writeMux.Lock()
	defer c.writeMux.Unlock()
	var body *anypb.Any
	body, err := anypb.New(msg)
	if err != nil {
		return 0, err
	}
	var m = packet.Message{
		Body: body,
	}
	data, err := proto.Marshal(&m)
	if err != nil {
		return 0, err
	}
	p := &packet.Packet{
		MessageType: msgType,
		CodecType:   packet.CodecType_Protobuf,
		Data:        data,
	}
	return c.codec.Encode(c.netConn, p)
}

func (c *Conn) Read(b []byte) (n int, err error) {
	for {
		switch {
		// case c.closeSignal != nil:
		// 	return 0, io.EOF
		case c.dataBuf.Len() > 0:
			return c.dataBuf.Read(b)
		default:
			c.cond.Wait()
			continue
		}
	}
}

func (c *Conn) Write(b []byte) (n int, err error) {
	c.writeMux.Lock()
	defer c.writeMux.Unlock()
	p := &packet.Packet{
		MessageType: packet.MessageType_Stream_Data,
		Data:        b,
	}

	return c.codec.Encode(c.netConn, p)
}

func (c *Conn) Close() error {
	c.closeOnce.Do(func() {
		if c.netConn != nil {
			c.netConn.Close()
		}
		close(c.closeSignal)
	})
	return nil
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

func (c *Conn) Stop() {
	c.netConn = nil
	c.Close()
}

func (c *Conn) GetConn() net.Conn {
	return c.netConn
}
