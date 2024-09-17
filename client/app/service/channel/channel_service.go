package channel

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"noah/client/app/entitie"
	"noah/client/app/service"
	"os"
	"os/exec"
	"runtime"
	"sync"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"

	ws "noah/client/app/infrastructure/websocket"
)

type Service struct {
	channels map[string]*Channel
}

type Channel struct {
	channelId   string
	channelType ChannelType
	conn        *websocket.Conn
	pty         *os.File
	tcpConn     net.Conn
	closeSig    chan struct{}
	mu          sync.Mutex // 添加互斥锁
	readFunc    func() error
	WriteFunc   func(msgType int, data []byte) error
}

type windowSize struct {
	High  int `json:"high"`
	Width int `json:"width"`
}

type ChannelType int

const (
	UnknownChannelType ChannelType = iota
	Pty
	Tcp
)

const (
	// Maximum message size allowed from peer.
	maxMessageSize = 5120
)

func NewService() service.Channel {
	return &Service{
		channels: make(map[string]*Channel),
	}
}

func (t *Service) NewChannel(channelId string, cType int, wsc *websocket.Conn, addr string) error {
	_, ok := t.channels[channelId]
	if ok {
		log.Println("Client already exists: ", channelId)
		return nil
	}

	channelType := ChannelType(cType)

	ch := &Channel{
		channelId:   channelId,
		channelType: channelType,
		conn:        wsc,
		closeSig:    make(chan struct{}),
		mu:          sync.Mutex{},
	}

	if channelType == Pty {
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case `windows`:
			cmd = exec.Command("cmd")
		case `linux`:
			cmd = exec.Command("bash")
		case `darwin`:
			cmd = exec.Command("zsh")
		default:
			return service.ErrUnsupportedPlatform
		}

		// 打开 pty
		ptmx, err := pty.Start(cmd)
		if err != nil {
			log.Println("Error starting PTY:", err)
			return fmt.Errorf("failed to start PTY: %w", err)
		}

		ch.pty = ptmx
		ch.readFunc = ch.ptyRead
		ch.WriteFunc = ch.ptyWrite

	}

	if channelType == Tcp {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return err
		}

		ch.tcpConn = conn
		ch.readFunc = ch.tcpRead
		ch.WriteFunc = ch.tcpWrite
	}

	t.channels[channelId] = ch

	go ch.readFunc()

	return nil
}

func (t *Service) Write(wsMessageType int, channelId string, data []byte) error {
	fmt.Println(t.channels)
	ch, ok := t.channels[channelId]
	if !ok {
		return errors.New("channel not found")
	}
	fmt.Println("Write to channel: ", channelId)

	return ch.WriteFunc(wsMessageType, data)
}

func (c *Channel) ptyRead() error {
	defer c.close()

	data := make([]byte, maxMessageSize)

	for {
		n, readErr := c.pty.Read(data)
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				break
			}
		}

		if n > 0 {
			ws_data := entitie.ChannelRequest{
				Action:      "write",
				ChannelId:   c.channelId,
				ChannelData: data[:n],
			}
			if err := ws.WriteMessage(c.conn, 0, entitie.MessageTypeChannel, ws_data, ""); err != nil {
				break
			}
		}
	}

	return nil
}

func (c *Channel) tcpRead() error {
	defer c.close()

	buffer := make([]byte, 1024)
	for {
		n, err := c.tcpConn.Read(buffer)
		if err != nil {
			log.Printf("Error reading from TCP: %v", err)
			break
		}

		ws_data := entitie.ChannelRequest{
			Action:      "write",
			ChannelId:   c.channelId,
			ChannelData: buffer[:n],
		}
		if err := ws.WriteMessage(c.conn, 0, entitie.MessageTypeChannel, ws_data, ""); err != nil {
			break
		}
	}

	return nil
}

func (c *Channel) ptyWrite(msgType int, data []byte) error {
	if msgType != websocket.BinaryMessage {
		if _, err := c.pty.Write(data); err != nil {
			return err
		}
	}

	var wdSize windowSize
	if err := json.Unmarshal(data, &wdSize); err != nil {
		return err
	}

	if err := pty.Setsize(c.pty, &pty.Winsize{
		Rows: uint16(wdSize.High),
		Cols: uint16(wdSize.Width),
	}); err != nil {
		return err
	}

	return nil
}

func (c *Channel) tcpWrite(msgType int, data []byte) error {
	_, err := c.tcpConn.Write(data)
	if err != nil {
		log.Println("Error writing to TCP:", err)
		return err
	}
	return nil
}

func (c *Channel) close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.pty != nil {
		log.Println("[!] Closing pty")
		if err := c.pty.Close(); err != nil {
			return err
		}
		c.pty = nil
		log.Println("[!] Closing pty end")
	}

	return nil
}
