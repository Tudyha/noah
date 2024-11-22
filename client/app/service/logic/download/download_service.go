package download

import (
	"fmt"
	"net/http"
	"noah/client/app/gateway"
	"noah/pkg/utils"

	"github.com/samber/do/v2"
)

type downloadService struct {
	gateway gateway.Gateway
}

func NewDownloadService(i do.Injector) (downloadService, error) {
	return downloadService{
		gateway: do.MustInvoke[gateway.Gateway](i),
	}, nil
}

func (d downloadService) DownloadFile(filename string, filepath string) error {
	url := fmt.Sprintf("/file/download/%s", filename)

	res, err := d.gateway.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	if err := utils.WriteFile(filepath, res); err != nil {
		return err
	}
	return nil
}
