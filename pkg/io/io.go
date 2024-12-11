package io

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

type WebSocketReaderWriterCloser struct {
	Conn *websocket.Conn
}

func (w *WebSocketReaderWriterCloser) Read(p []byte) (n int, err error) {
	_, reader, err := w.Conn.NextReader()
	if err != nil {
		return
	}
	return reader.Read(p)
}

func (w *WebSocketReaderWriterCloser) Write(p []byte) (n int, err error) {
	return len(p), w.Conn.WriteMessage(websocket.TextMessage, p)
}

func (w *WebSocketReaderWriterCloser) Close() error {
	return w.Conn.Close()
}

type PtyReaderWriterCloser struct {
	IO *os.File
}

func (w *PtyReaderWriterCloser) Read(p []byte) (n int, err error) {
	return w.IO.Read(p)
}

type ptyData struct {
	Type  string `json:"type"`
	Data  any    `json:"data"`
	High  int    `json:"high"`
	Width int    `json:"width"`
}

func (w *PtyReaderWriterCloser) Write(p []byte) (n int, err error) {
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

func (w *PtyReaderWriterCloser) Close() error {
	return w.IO.Close()
}

func Copy(left, right io.ReadWriter) error {
	var err, err1 error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err1 = io.Copy(right, left)
	}()
	_, err = io.Copy(left, right)
	wg.Wait()
	if err1 != nil && !errors.Is(err1, os.ErrDeadlineExceeded) { // requires Go 1.15+
		return err1
	}
	if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
		return err
	}
	return nil
}
