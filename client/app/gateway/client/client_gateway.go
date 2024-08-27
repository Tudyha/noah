package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"noah/client/app/environment"
	"noah/client/app/gateway"
)

type Gateway struct {
	Configuration *environment.Configuration
	HttpClient    *http.Client
}

func NewGateway(configuration *environment.Configuration, httpClient *http.Client) gateway.Gateway {
	return &Gateway{
		Configuration: configuration,
		HttpClient:    httpClient,
	}
}

func (c Gateway) NewRequest(method string, url string, body []byte) (*gateway.HttpResponse, error) {
	req, err := http.NewRequest(method, fmt.Sprint(c.Configuration.Server.Url, url), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", c.Configuration.Connection.Token)

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed with status code %d", res.StatusCode)
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response gateway.HttpResponse

	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c Gateway) NewFileDownloadRequest(method string, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, fmt.Sprint(c.Configuration.Server.Url, url), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", c.Configuration.Connection.Token)

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed with status code %d", res.StatusCode)
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}
