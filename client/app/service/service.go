package service

import (
	"errors"
	"noah/client/app/entitie"

	"github.com/gorilla/websocket"
)

var (
	ErrUnsupportedPlatform = errors.New("unsupported platform")
	ErrDeadlineExceeded    = errors.New("command deadline exceeded")
)

type Services struct {
	Information
	Terminal
	Pty
	Download
}

type Information interface {
	LoadDeviceSpecs() (*entitie.Device, error)
}

type Terminal interface {
	Run(command string) ([]byte, error)
}

type Pty interface {
	Run(wsc *websocket.Conn) error
}

type Download interface {
	DownloadFile(filename string, filepath string) error
}
