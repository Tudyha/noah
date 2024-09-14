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
	Command
	Pty
	Download
	FileExplorer
}

type Information interface {
	LoadClientSpecs() (*entitie.Client, error)
	GetSystemInfo() (*entitie.SystemInfo, error)
}

type Command interface {
	Run(command string) ([]byte, error)
	GetProcessList() ([]entitie.Process, error)
	KillProcess(pid int32) error
}

type Pty interface {
	Run(wsc *websocket.Conn) error
	NewPtyClient(channelId string, wsc *websocket.Conn) error
	Write(msgType int, channelId string, data []byte) error
	SetSize(channelId string, data []byte) error
}

type Download interface {
	DownloadFile(filename string, filepath string) error
}

type FileExplorer interface {
	GetFileExplorer(path string) ([]entitie.FileExplorer, error)
	ReadFile(path string) ([]byte, error)
	Rename(path string, newFilename string) error
	Remove(path string) error
	WriteFile(path string, content []byte) error
	MkDir(path string) error
}
