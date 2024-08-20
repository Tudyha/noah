package pty

import (
	"encoding/json"
	"fmt"
	"io"
	"noah/client/app/service"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

type Service struct {
}

type ptyClient struct {
	conn      *websocket.Conn
	pty       *os.File
	closeSig  chan struct{}
	readDone  chan struct{}
	writeDone chan struct{}
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

	//打开pty
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return fmt.Errorf("failed to start PTY: %w", err)
	}

	ptyClient := &ptyClient{
		conn:      wsc,
		pty:       ptmx,
		closeSig:  make(chan struct{}, 1),
		readDone:  make(chan struct{}),
		writeDone: make(chan struct{}),
	}
	defer func() {
		ptyClient.Close()
		close(ptyClient.closeSig)
	}()

	go func() {
		err := ptyClient.read()
		if err != nil {
			ptyClient.Log("[!] Error ptyClient read: ", err.Error())
			close(ptyClient.readDone)
		}
	}()
	go func() {
		err := ptyClient.write()
		if err != nil {
			ptyClient.Log("[!] Error ptyClient write: ", err.Error())
			close(ptyClient.writeDone)
		}
	}()

	<-ptyClient.closeSig

	return nil
}

func (c *ptyClient) Log(v ...any) {
	fmt.Println(v...)
}

func (c *ptyClient) write() error {
	defer func() {
		c.closeSig <- struct{}{}
	}()

	data := make([]byte, maxMessageSize)

	for {
		//time.Sleep(10 * time.Millisecond)
		n, readErr := c.pty.Read(data)
		if readErr != nil {
			return readErr
		}
		//if n == 0 {
		//	break
		//}
		if n > 0 {
			if err := c.conn.WriteMessage(websocket.TextMessage, data[:n]); err != nil {
				return err
			}
		}
	}
}

func (c *ptyClient) read() error {
	defer func() {
		c.closeSig <- struct{}{}
	}()

	var zeroTime time.Time
	c.conn.SetReadDeadline(zeroTime)

	for {
		msgType, connReader, err := c.conn.NextReader()
		if err != nil {
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
}

func (c *ptyClient) Close() error {
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
