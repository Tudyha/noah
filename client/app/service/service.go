package service

import (
	"errors"
	"noah/client/app/entitie"
)

var (
	ErrUnsupportedPlatform = errors.New("unsupported platform")
	ErrDeadlineExceeded    = errors.New("command deadline exceeded")
)

type Services struct {
	Information
	Command
	Download
	FileExplorer
}

type Information interface {
	LoadClientSpecs() (*entitie.Client, error)
	GetSystemInfo() (*entitie.SystemInfo, error)
	GetProcessList() ([]entitie.Process, error)
	KillProcess(pid int32) error
	GetNetworkInfo() ([]entitie.NetworkInfo, error)
	GetDockerContainerList() (res []entitie.DockerContainer, err error)
}

type Command interface {
	Run(command string) ([]byte, error)
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
