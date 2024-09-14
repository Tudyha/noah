package pty

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"noah/client/app/entitie"
	"noah/client/app/service"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"

	ws "noah/client/app/infrastructure/websocket"
)

type Service struct{}

var (
	ptyClients = make(map[string]*ptyClient)
)

type ptyClient struct {
	conn     *websocket.Conn
	pty      *os.File
	closeSig chan struct{}
	mu       sync.Mutex // 添加互斥锁
}

type windowSize struct {
	High  int `json:"high"`
	Width int `json:"width"`
}

const (
	// Maximum message size allowed from peer.
	maxMessageSize = 5120
)

func NewService() service.Pty {
	return &Service{}
}

func (t *Service) Run(wsc *websocket.Conn) error {
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
		return fmt.Errorf("failed to start PTY: %w", err)
	}

	ptyClient := &ptyClient{
		conn:     wsc,
		pty:      ptmx,
		closeSig: make(chan struct{}),
	}
	defer func() {
		if err := ptyClient.close(); err != nil {
			ptyClient.Log("[!] Error closing ptyClient: ", err)
		}
		close(ptyClient.closeSig)
	}()

	go func() {
		err := ptyClient.read()
		if err != nil {
			ptyClient.Log("[!] Error ptyClient read: ", err.Error())
		}
	}()
	go func() {
		err := ptyClient.write()
		if err != nil {
			ptyClient.Log("[!] Error ptyClient write: ", err.Error())
		}
	}()

	<-ptyClient.closeSig

	return nil
}

func (c *ptyClient) Log(v ...any) {
	log.Println(v...)
}

func (c *ptyClient) write() error {
	defer func() {
		if err := c.close(); err != nil {
			c.Log("[!] Error closing ptyClient: ", err)
		}
		c.closeSig <- struct{}{}
	}()

	data := make([]byte, maxMessageSize)

	for {
		n, readErr := c.pty.Read(data)
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				break
			}
			return readErr
		}

		if n > 0 {
			if err := c.conn.WriteMessage(websocket.TextMessage, data[:n]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *ptyClient) read() error {
	defer func() {
		if err := c.close(); err != nil {
			c.Log("[!] Error closing ptyClient: ", err)
		}
		c.closeSig <- struct{}{}
	}()

	var zeroTime time.Time
	c.conn.SetReadDeadline(zeroTime)

	for {
		msgType, connReader, err := c.conn.NextReader()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		if msgType != websocket.BinaryMessage {
			if _, err := io.Copy(c.pty, connReader); err != nil {
				return err
			}
			continue
		}

		data := make([]byte, maxMessageSize)
		n, err := connReader.Read(data)
		if err != nil {
			return err
		}

		var wdSize windowSize
		if err := json.Unmarshal(data[:n], &wdSize); err != nil {
			return err
		}

		if err := pty.Setsize(c.pty, &pty.Winsize{
			Rows: uint16(wdSize.High),
			Cols: uint16(wdSize.Width),
		}); err != nil {
			return err
		}
	}
	return nil
}

func (c *ptyClient) close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}
		c.conn = nil
	}

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

func (t *Service) NewPtyClient(channelId string, wsc *websocket.Conn) error {
	_, ok := ptyClients[channelId]
	if ok {
		log.Println("Client already exists: ", channelId)
		return nil
	}

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
		return fmt.Errorf("failed to start PTY: %w", err)
	}

	cl := &ptyClient{
		conn: wsc,
		pty:  ptmx,
	}

	ptyClients[channelId] = cl

	go cl.read1(channelId)

	return nil
}

func (c *ptyClient) read1(channelId string) {
	data := make([]byte, maxMessageSize)

	for {
		n, readErr := c.pty.Read(data)
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				c.pty.Close()
				break
			}
		}

		if n > 0 {
			ws_data := entitie.PtyRequest{
				Action:      "write",
				ChannelId:   channelId,
				ChannelData: data[:n],
			}
			if err := ws.WriteMessage(c.conn, 0, entitie.MessageTypePty, ws_data, ""); err != nil {
				c.pty.Close()
				break
			}
		}
	}
}

func (t *Service) Write(msgType int, channelId string, data []byte) error {
	c, ok := ptyClients[channelId]
	if !ok {
		return errors.New("client not found")
	}

	if msgType != websocket.BinaryMessage {
		if _, err := c.pty.Write(data); err != nil {
			return err
		}
	}

	t.SetSize(channelId, data)

	return nil
}

func (t *Service) SetSize(channelId string, data []byte) error {
	c, ok := ptyClients[channelId]
	if !ok {
		return errors.New("client not found")
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
