package pty

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"noah/client/app/service"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

type Service struct{}

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
		if err := c.pty.Close(); err != nil {
			return err
		}
		c.pty = nil
	}

	return nil
}
