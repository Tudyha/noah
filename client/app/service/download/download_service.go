package download

import (
	"fmt"
	"net/http"
	"noah/client/app/gateway"
	"noah/client/app/service"
	"noah/client/app/utils"
)

type Service struct {
	Gateway gateway.Gateway
}

func NewService(gateway gateway.Gateway) service.Download {
	return &Service{
		Gateway: gateway,
	}
}

func (d Service) DownloadFile(filename string, filepath string) error {
	url := fmt.Sprintf("/download/%s", filename)

	res, err := d.Gateway.NewFileDownloadRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	if err := utils.WriteFile(filepath, res); err != nil {
		return err
	}
	return nil
}
