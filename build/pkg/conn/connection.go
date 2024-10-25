package conn

import (
	"bytes"
	"io"
	"net"
	"sync"
	"time"
)

type Conn struct {
	net.Conn
	connId         uint32
	mux            *Mux
	connStatusOkCh chan struct{}
	receiveChan    chan []byte
	once           sync.Once
	lk             LinkInfo
}

func NewConn(id uint32, m *Mux) *Conn {
	return &Conn{
		connId:         id,
		mux:            m,
		connStatusOkCh: make(chan struct{}),
		receiveChan:    make(chan []byte, 32),
		once:           sync.Once{},
	}
}

func (c *Conn) receive(data []byte) {
	c.receiveChan <- data
}

func (c *Conn) Read(b []byte) (n int, err error) {
	data, ok := <-c.receiveChan
	if !ok {
		return 0, io.EOF
	}
	copy(b, data)
	return len(data), nil
}

func (c *Conn) Write(b []byte) (n int, err error) {
	n = len(b)
	d := make([]byte, n)
	copy(d[:], b[:n])
	c.mux.writeMsg(data, c.connId, d)
	return n, nil
}

func (c *Conn) Close() error {
	c.once.Do(func() {
		c.mux.removeConn(c.connId)
		close(c.receiveChan)
		c.mux.writeMsg(connClose, c.connId, nil)
	})
	return nil
}

// Copy fixme 加解密、压缩
func (c *Conn) Copy(target io.ReadWriter) {
	go func() {
		io.Copy(target, c)
	}()
	io.Copy(c, target)
}

func (c *Conn) ReadFull() ([]byte, error) {
	var buffer bytes.Buffer
	buf := make([]byte, 1024*1024)
	for {
		n, err := c.Read(buf)
		if err != nil {
			break
		}
		if n > 0 {
			buffer.Write(buf[:n])
		}
	}
	return buffer.Bytes(), nil
}

func (c *Conn) GetLk() LinkInfo {
	return c.lk
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
