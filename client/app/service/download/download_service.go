package download

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"noah/client/app/gateway"
	"noah/client/app/service"
	"os"
	"path/filepath"
)

type Service struct {
	Gateway gateway.Gateway
}

func NewService(gateway gateway.Gateway) service.Download {
	return &Service{
		Gateway: gateway,
	}
}

func (d Service) DownloadFile(filename string, filepath string) ([]byte, error) {
	url := fmt.Sprintf("/download/%s", filename)

	res, err := d.Gateway.NewFileDownloadRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(filepath, res, os.ModePerm); err != nil {
		return nil, err
	}
	return []byte(filename), nil
}

func getFilenameFromPath(path string) string {
	return filepath.Base(path)
}
