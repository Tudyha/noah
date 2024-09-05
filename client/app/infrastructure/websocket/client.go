package websocket

import (
	"fmt"
	"net/http"
	"net/url"
	"noah/client/app/environment"
	"strings"

	"github.com/gorilla/websocket"
)

func NewConnection(configuration *environment.Configuration, path string) (*websocket.Conn, error) {
	host := configuration.Server.Address
	host = strings.TrimPrefix(host, "http://")
	host = strings.TrimPrefix(host, "https://")
	host = strings.TrimSuffix(host, "/")

	if configuration.Server.HttpPort != "" {
		host = fmt.Sprint(host, ":", configuration.Server.HttpPort)
	}

	scheme := "ws"
	if strings.Contains(configuration.Server.Address, "https") {
		scheme = "wss"
	}

	u := url.URL{Scheme: scheme, Host: host, Path: path}

	header := http.Header{}
	// header.Set("x-client", clientID)
	header.Set("Authorization", configuration.Connection.Token)

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	return conn, err
}
