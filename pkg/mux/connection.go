package mux

import (
	"bytes"
	"io"
	"net"
	"sync"
	"time"
)

type Conn struct {
	net.Conn
	id             uint32
	mux            *Mux
	connStatusOkCh chan struct{}
	receiveChan    chan []byte
	once           sync.Once
	network        string
	addr           string
	internalBuffer bytes.Buffer
}

func NewConn(id uint32, m *Mux) *Conn {
	return &Conn{
		id:             id,
		mux:            m,
		connStatusOkCh: make(chan struct{}),
		receiveChan:    make(chan []byte, 32),
		once:           sync.Once{},
		internalBuffer: bytes.Buffer{},
	}
}

func (c *Conn) receive(data []byte) {
	c.receiveChan <- data
}

func (c *Conn) Read(b []byte) (n int, err error) {
	// 先从内部缓冲区读取数据
	if c.internalBuffer.Len() > 0 {
		n = copy(b, c.internalBuffer.Bytes())
		c.internalBuffer.Next(n)
		if n == len(b) {
			return n, nil
		}
	}

	// 如果内部缓冲区为空或不足以填满 b，从 receiveChan 读取数据
	data, ok := <-c.receiveChan
	if !ok {
		return n, io.EOF
	}

	// 计算还需要读取多少数据
	remaining := len(b) - n
	if remaining > 0 {
		n += copy(b[n:], data)
		if len(data) > remaining {
			c.internalBuffer.Write(data[remaining:])
		}
	} else {
		// 如果 b 已经满了，将剩余的数据存储到内部缓冲区
		c.internalBuffer.Write(data)
	}

	return n, nil
}

func (c *Conn) Write(b []byte) (n int, err error) {
	n = len(b)
	d := make([]byte, n)
	copy(d[:], b[:n])
	c.mux.writeMsg(flagData, c.id, d)
	return n, nil
}

func (c *Conn) Close() error {
	c.once.Do(func() {
		c.mux.removeConn(c.id)
		close(c.receiveChan)
	})
	return nil
}

func (c *Conn) ReadFull() ([]byte, error) {
	var buffer bytes.Buffer
	buf := make([]byte, 32*1024)
	for {
		n, err := c.Read(buf)
		if n > 0 {
			buffer.Write(buf[:n])
		}
		if err != nil {
			break
		}
	}
	return buffer.Bytes(), nil
}

func (c *Conn) LocalAddr() net.Addr {
	return c.mux.Addr()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.mux.RemoteAddr()
}

func (c *Conn) SetDeadline(t time.Time) error {
	return nil
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (c *Conn) GetNetwork() string {
	return c.network
}

func (c *Conn) GetAddr() string {
	return c.addr
}
