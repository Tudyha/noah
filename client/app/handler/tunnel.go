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

type TunnelHandler struct{}

func NewTunnelHandler() conn.MessageHandler {
	return &TunnelHandler{}
}

func (p *TunnelHandler) Handle(ctx conn.Context) error {
	var err error
	defer func() {
		if err != nil {
			ctx.GetConn().WriteProtoMessage(packet.MessageType_Tunnel_Open_Ack, &packet.OpenTunnelAck{
				Code: -1,
				Msg:  err.Error(),
			})
		}
	}()
	var msg packet.OpenTunnel
	if err = ctx.ShouldBindProto(&msg); err != nil {
		return err
	}

	var target io.ReadWriteCloser
	switch msg.TunnelType {
	case packet.OpenTunnel_PTY:
		target, err = p.openPty()
		if err != nil {
			return err
		}
	default:
		err = fmt.Errorf("unsupported tunnel type: %s", msg.TunnelType)
		return err
	}

	if err := ctx.GetConn().WriteProtoMessage(packet.MessageType_Tunnel_Open_Ack, &packet.OpenTunnelAck{
		Code: 0,
		Msg:  "",
	}); err != nil {
		return err
	}

	// 劫持底层连接
	src, err := ctx.Hijack()
	if err != nil {
		return err
	}

	go func() {
		_, err := io.Copy(src, target)
		if err != nil {
			target.Close()
			src.Close()
			return
		}
	}()

	go func() {
		_, err := io.Copy(target, src)
		if err != nil {
			target.Close()
			src.Close()
			return
		}
	}()

	return nil
}

func (p *TunnelHandler) MessageType() packet.MessageType {
	return packet.MessageType_Tunnel_Open
}

func (h *TunnelHandler) openPty() (io.ReadWriteCloser, error) {
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
