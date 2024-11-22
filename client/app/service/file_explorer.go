package service

import "noah/pkg/mux/message"

type IFileExplorerService interface {
	GetFileExplorer(path string) ([]message.FileExplorer, error)
	ReadFile(path string) ([]byte, error)
	Rename(path string, newFilename string) error
	Remove(path string) error
	WriteFile(path string, content []byte) error
	MkDir(path string) error
}
