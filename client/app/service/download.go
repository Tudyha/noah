package service

type IDownloadService interface {
	DownloadFile(filename string, filepath string) error
}
