package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"noah/pkg/conn"
	"noah/pkg/packet"
	"os"
	"os/exec"
	"runtime"

	"github.com/creack/pty"
)

type PtyHandler struct{}

func NewPtyHandler() conn.MessageHandler {
	return &PtyHandler{}
}

func (p *PtyHandler) Handle(ctx conn.Context) error {
	rwc, err := p.open()
	if err != nil {
		return err
	}
	defer rwc.Close()

	conn, err := ctx.GetConn().GetSmuxSession().Open()
	if err != nil {
		return err
	}
	defer conn.Close()

	go func() {
		io.Copy(conn, rwc)
	}()

	go func() {
		io.Copy(rwc, conn)
	}()

	return nil
}

func (p *PtyHandler) MessageType() packet.MessageType {
	return packet.MessageType_Tunnel_Pty
}

func (h *PtyHandler) open() (io.ReadWriteCloser, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case `linux`:
		cmd = exec.Command("bash")
	case `darwin`:
		cmd = exec.Command("zsh")
	default:
		return nil, fmt.Errorf("unsupported os: %s", runtime.GOOS)
	}

	// 打开 pty
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}
	return &ptyReaderWriterCloser{IO: ptmx}, nil
}

type ptyReaderWriterCloser struct {
	IO *os.File
}

func (w *ptyReaderWriterCloser) Read(p []byte) (n int, err error) {
	return w.IO.Read(p)
}

type ptyData struct {
	Type  string `json:"type"`
	Data  any    `json:"data"`
	High  int    `json:"high"`
	Width int    `json:"width"`
}

func (w *ptyReaderWriterCloser) Write(p []byte) (n int, err error) {
	n = len(p)
	var ptyData ptyData
	err = json.Unmarshal(p, &ptyData)
	if err != nil {
		return
	}
	switch ptyData.Type {
	case "size":
		if err = pty.Setsize(w.IO, &pty.Winsize{
			Rows: uint16(ptyData.High),
			Cols: uint16(ptyData.Width),
		}); err != nil {
			return
		}
	case "data":
		_, err = w.IO.Write([]byte(ptyData.Data.(string)))
	}
	return
}

func (w *ptyReaderWriterCloser) Close() error {
	return w.IO.Close()
}
