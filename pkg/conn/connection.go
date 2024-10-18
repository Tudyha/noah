package conn

import (
	"net"
	"sync"
	"time"
)

type Conn struct {
	net.Conn
	connId         uint32
	mux            *Mux
	connStatusOkCh chan struct{}
	recieveChan    chan []byte
	once           sync.Once
	lk             LinkInfo
}

func NewConn(id uint32, m *Mux) *Conn {
	return &Conn{
		connId:         id,
		mux:            m,
		connStatusOkCh: make(chan struct{}),
		recieveChan:    make(chan []byte, 32),
		once:           sync.Once{},
	}
}

func (c *Conn) recieve(data []byte) {
	c.recieveChan <- data
}

func (c *Conn) Read(b []byte) (n int, err error) {
	data := <-c.recieveChan
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
		close(c.recieveChan)
	})
	return nil
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
