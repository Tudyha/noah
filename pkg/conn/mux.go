package conn

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type Mux struct {
	net.Listener
	conn             *websocket.Conn  // websocket连接
	waitConnQueue    chan *Conn       // 等待处理的连接
	conns            map[uint32]*Conn // 连接池
	connId           uint32           // 自增连接id
	sendMessageQueue chan Message     // 发送消息队列
	once             sync.Once
	Closed           bool
}

// NewMux 基于websocket连接的多路复用器
func NewMux(c *websocket.Conn) *Mux {
	m := &Mux{
		conn:             c,
		waitConnQueue:    make(chan *Conn),
		conns:            make(map[uint32]*Conn, 32),
		sendMessageQueue: make(chan Message, 32),
		once:             sync.Once{},
	}

	go m.read()
	go m.write()

	return m
}

// Accept new connection.
func (m *Mux) Accept() (net.Conn, error) {
	c, ok := <-m.waitConnQueue
	if !ok {
		return nil, errors.New("accept fail")
	}
	return c, nil
}

// close mux
func (m *Mux) Close() error {
	m.once.Do(func() {
		m.conn.Close()
		m.Closed = true
		// close(m.waitConnQueue)
	})
	return nil
}

// Addr returns local address.
func (m *Mux) Addr() net.Addr {
	return m.conn.LocalAddr()
}

// RemoteAddr returns remote address.
func (m *Mux) RemoteAddr() net.Addr {
	return m.conn.RemoteAddr()
}

// read websocket message.
func (m *Mux) read() {
	defer m.Close()

	for {
		_, b, err := m.conn.ReadMessage()
		if err != nil {
			// if err == io.EOF {
			// 	// log
			// } else {
			// 	// log
			// }
			return
		}

		var message Message

		err = json.Unmarshal(b, &message)
		if err != nil {
			continue
		}
		switch message.Flag {
		case newConn:
			var lk LinkInfo
			err := json.Unmarshal(message.Data, &lk)
			if err != nil {
				continue
			}
			c := NewConn(message.ConnId, m)
			c.lk = lk
			m.conns[c.connId] = c
			m.waitConnQueue <- c
			m.writeMsg(newConnOk, message.ConnId, nil)
		case newConnOk:
			if conn, ok := m.conns[message.ConnId]; ok {
				conn.connStatusOkCh <- struct{}{}
			}
		case data:
			if conn, ok := m.conns[message.ConnId]; ok {
				conn.recieve(message.Data)
			}
		}

	}
}

func (m *Mux) write() {
	defer m.Close()

	for message := range m.sendMessageQueue {
		b, _ := json.Marshal(message)
		if err := m.conn.WriteMessage(websocket.BinaryMessage, b); err != nil {
			return
		}
	}
}

func (m *Mux) NewConn(network string, addr string) (*Conn, error) {
	conn := NewConn(m.getConnId(), m)
	//it must be Set before send
	m.conns[conn.connId] = conn
	lk := LinkInfo{
		Addr:    addr,
		Network: network,
	}
	data, err := json.Marshal(lk)
	if err != nil {
		return nil, err
	}
	m.writeMsg(newConn, conn.connId, data)
	//Set a timer timeout 120 second
	timer := time.NewTimer(time.Minute * 2)
	defer timer.Stop()
	select {
	case <-conn.connStatusOkCh:
		return conn, nil
	case <-timer.C:
	}
	return nil, errors.New("create connection fail, the server refused the connection")
}

func (m *Mux) writeMsg(flag flag, connId uint32, b []byte) {
	m.sendMessageQueue <- Message{
		Flag:   flag,
		ConnId: connId,
		Data:   b,
	}
}

func (m *Mux) getConnId() (id uint32) {
	id = atomic.AddUint32(&m.connId, 1)
	return
}

func (m *Mux) removeConn(id uint32) {
	delete(m.conns, id)
}

func (m *Mux) healthCheck() {
	for {
		time.Sleep(time.Second * 5)
		fmt.Println("mux health check", m.conns)
	}
}
