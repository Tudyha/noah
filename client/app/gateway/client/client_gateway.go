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
	req.Header.Set("Token", c.Configuration.Connection.Token)

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		if res.StatusCode == 401 {
			//尝试刷新token
			if err := c.refreshToken(); err != nil {
				return nil, err
			}
			//重新请求
			return c.NewRequest(method, url, body)
		} else {
			return nil, fmt.Errorf("failed with status code %d", res.StatusCode)
		}
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

func (c Gateway) refreshToken() error {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprint(c.Configuration.Server.Url, "/refresh_token?refreshToken=", c.Configuration.Connection.RefreshToken), nil)
	if err != nil {
		return err
	}

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("failed refresh token")
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var response gateway.HttpResponse

	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return err
	}

	if response.Code != 0 {
		return fmt.Errorf("failed refresh token")
	}

	if data, ok := response.Data.(map[string]interface{}); ok {
		if data["token"] == nil || data["refreshToken"] == nil {
			return fmt.Errorf("failed refresh token")
		}
		c.Configuration.Connection.Token = data["token"].(string)
		c.Configuration.Connection.RefreshToken = data["refreshToken"].(string)
	} else {
		return fmt.Errorf("failed refresh token")
	}

	return nil
}

func (c Gateway) NewFileDownloadRequest(method string, url string, body []byte) ([]byte, error) {
	if method == "" || url == "" {
		return nil, fmt.Errorf("method or URL cannot be empty")
	}
	req, err := http.NewRequest(method, fmt.Sprint(c.Configuration.Server.Url, url), bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return nil, err
	}
	req.Header.Set("Token", c.Configuration.Connection.Token)

	c.HttpClient.Timeout = 0

	res, err := c.HttpClient.Do(req)
	if err != nil {
		fmt.Printf("Error executing request: %v\n", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed with status code %d", res.StatusCode)
	}
	// 使用较大的缓冲区来提高读取速度
	const bufferSize = 4096 // 可以根据实际情况调整这个值
	buffer := make([]byte, bufferSize)
	bodyBytes := &bytes.Buffer{}

	_, err = io.CopyBuffer(bodyBytes, res.Body, buffer)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return nil, err
	}

	return bodyBytes.Bytes(), nil
}
