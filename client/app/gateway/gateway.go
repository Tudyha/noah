package gateway

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"noah/client/app/environment"

	"github.com/samber/do/v2"
)

type Gateway struct {
	Env        *environment.Environment
	HttpClient http.Client
}

func NewGateway(i do.Injector) (Gateway, error) {
	return Gateway{
		Env:        do.MustInvoke[*environment.Environment](i),
		HttpClient: http.Client{},
	}, nil
}

func (c Gateway) NewRequest(method string, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Token", c.Configuration.Connection.Token)

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		// if res.StatusCode == 401 {
		// 	//尝试刷新token
		// 	if err := c.refreshToken(); err != nil {
		// 		return nil, err
		// 	}
		// 	//重新请求
		// 	return c.NewRequest(method, url, body)
		// } else {
		// 	return nil, fmt.Errorf("failed with status code %d", res.StatusCode)
		// }
		return nil, fmt.Errorf("failed with status code %d", res.StatusCode)
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}
