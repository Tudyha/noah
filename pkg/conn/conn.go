package conn

import (
	"bytes"
	"net"
	"noah/pkg/packet"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Conn struct {
	netConn   net.Conn     // 底层连接
	codec     packet.Codec // 编解码器
	closeOnce sync.Once

	dataBuf  *bytes.Buffer // 数据
	writeMux *sync.Mutex
}

func NewConn(conn net.Conn) *Conn {
	c := &Conn{
		codec:    packet.NewCodec(),
		netConn:  conn,
		dataBuf:  bytes.NewBuffer(nil),
		writeMux: &sync.Mutex{},
	}

	return c
}

func (c *Conn) ReadMessage() (*packet.Packet, error) {
	return c.codec.Decode(c.netConn)
}

func (c *Conn) WriteProtoMessage(msgType packet.MessageType, msg proto.Message) error {
	c.writeMux.Lock()
	defer c.writeMux.Unlock()
	var body *anypb.Any
	body, err := anypb.New(msg)
	if err != nil {
		return err
	}
	var m = packet.Message{
		Body: body,
	}
	data, err := proto.Marshal(&m)
	if err != nil {
		return err
	}
	p := &packet.Packet{
		MessageType: msgType,
		CodecType:   packet.CodecType_Protobuf,
		Data:        data,
	}
	_, err = c.codec.Encode(c.netConn, p)
	return err
}

func (c *Conn) Read(b []byte) (n int, err error) {
	if c.dataBuf.Len() > 0 {
		return c.dataBuf.Read(b)
	}

	p, err := c.ReadMessage()
	if err != nil {
		return 0, err
	}

	c.dataBuf.Write(p.Data)
	return c.dataBuf.Read(b)
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
	})
	return nil
}

func (c *Conn) Release() {
	c.netConn = nil
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.netConn.RemoteAddr()
}

func (c *Conn) SetReadDeadline(t time.Time) {
	c.netConn.SetReadDeadline(t)
}

func (c *Conn) SetWriteDeadline(t time.Time) {
	c.netConn.SetWriteDeadline(t)
}

func (c *Conn) SetDeadline(t time.Time) {
	c.netConn.SetDeadline(t)
}
