package mux

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	"noah/pkg/mux/message"
)

type Mux struct {
	net.Listener
	reader        io.Reader
	writer        io.Writer
	connIdManager *iDManager // 连接id管理
	waitConnQueue chan *Conn // 等待处理的连接
	conns         *sync.Map  // 连接池

	sendPacketQueue chan packet // 发送消息队列

	closedCallbackHandler func() // mux关闭回调
	closeOnce             sync.Once

	pingHandler func() []byte     // 心跳连接ping，可在ping时附带数据
	pongHandler func(data []byte) // 心跳连接pong，可在pong时处理ping附带的数据
}

func NewMux(reader io.Reader, writer io.Writer) *Mux {
	m := &Mux{
		reader:                reader,
		writer:                writer,
		sendPacketQueue:       make(chan packet, 32),
		connIdManager:         newIDManager(),
		conns:                 &sync.Map{},
		waitConnQueue:         make(chan *Conn),
		closedCallbackHandler: nil,
		closeOnce:             sync.Once{},
		pingHandler:           nil,
		pongHandler:           nil,
	}

	return m
}

func (m *Mux) Start() error {
	go func() {
		defer m.Close()
		m.read()
	}()
	go func() {
		defer m.Close()
		m.write()
	}()
	go m.ping()

	return nil
}

// Accept new connection.
func (m *Mux) Accept() (net.Conn, error) {
	c, ok := <-m.waitConnQueue
	if !ok {
		return nil, errors.New("accept fail")
	}
	return c, nil
}

// Close mux
func (m *Mux) Close() error {
	m.closeOnce.Do(func() {
		close(m.waitConnQueue)
		if m.closedCallbackHandler != nil {
			m.closedCallbackHandler()
		}
	})
	return nil
}

func (m *Mux) SetClosedCallbackHandler(h func()) {
	m.closedCallbackHandler = h
}

func (m *Mux) SetPongHandler(h func(data []byte)) {
	m.pongHandler = h
}

func (m *Mux) SetPingHandler(h func() []byte) {
	m.pingHandler = h
}

// Addr returns local address.
func (m *Mux) Addr() net.Addr {
	return nil
}

// RemoteAddr returns remote address.
func (m *Mux) RemoteAddr() net.Addr {
	return nil
}

// read packet from reader.
func (m *Mux) read() error {
	buf := bufio.NewReader(m.reader)
	for {
		packet, err := readPacket(buf)
		if err != nil {
			return err
		}

		m.handlerPacket(packet)
	}
}

// handler packet
func (m *Mux) handlerPacket(packet *packet) {
	switch packet.Flag {
	case Flag_New_Conn:
		// new connection
		c := NewConn(packet.ConnId, m)

		var lk message.LinkInfo
		if err := json.Unmarshal(packet.Data, &lk); err != nil {
			return
		}

		c.network = lk.Network
		c.addr = lk.Addr

		m.conns.Store(c.id, c)
		m.waitConnQueue <- c
		m.writeMsg(Flag_Conn_Ok, packet.ConnId, nil)
	case Flag_Conn_Ok:
		// connection ok
		if conn, ok := m.conns.Load(packet.ConnId); ok {
			conn.(*Conn).connStatusOkCh <- struct{}{}
		}
	case Flag_Data:
		// conn data
		if conn, ok := m.conns.Load(packet.ConnId); ok {
			conn.(*Conn).receive(packet.Data)
		}
	case Flag_Close:
		// conn close
		if conn, ok := m.conns.Load(packet.ConnId); ok {
			conn.(*Conn).Close()
		}
	case Flag_Ping:
		// ping
		if m.pongHandler != nil {
			m.pongHandler(packet.Data)
		}
		m.writeMsg(Flag_Pong, 0, nil)
	}
}

func (m *Mux) write() error {
	for p := range m.sendPacketQueue {
		d, err := buildPacket(p.Flag, p.ConnId, p.Data)
		if err != nil {
			continue
		}
		n, err := m.writer.Write(d)
		if err != nil {
			return err
		}
		if n != len(d) {
			return err
		}
	}
	return nil
}

func (m *Mux) Dial(network, address string) (*Conn, error) {
	connId, err := m.connIdManager.GetID()
	if err != nil {
		return nil, err
	}
	conn := NewConn(connId, m)
	conn.network = network
	conn.addr = address
	// it must be Set before send
	m.conns.Store(connId, conn)

	data, err := json.Marshal(message.LinkInfo{
		Network: network,
		Addr:    address,
	})
	if err != nil {
		return nil, err
	}

	m.writeMsg(Flag_New_Conn, connId, data)
	// Set a timer timeout 120 second
	timer := time.NewTimer(time.Minute * 2)
	defer timer.Stop()
	select {
	case <-conn.connStatusOkCh:
		return conn, nil
	case <-timer.C:
	}
	return nil, errors.New("create connection fail, the server refused the connection")
}

func (m *Mux) writeMsg(flag flag, connId uint32, data []byte) {
	m.sendPacketQueue <- packet{
		Flag:   flag,
		Data:   data,
		ConnId: connId,
	}
}

func (m *Mux) removeConn(id uint32) {
	m.conns.Delete(id)
	m.connIdManager.ReleaseID(id)
	m.writeMsg(Flag_Close, id, nil)
}

func (m *Mux) ping() {
	for {
		var data []byte
		if m.pingHandler != nil {
			data = m.pingHandler()
		}
		m.writeMsg(Flag_Ping, 0, data)
		time.Sleep(time.Second * 30)
	}
}
