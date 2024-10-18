package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"noah/client/app/entitie"
	"noah/client/app/environment"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	mu sync.Mutex = sync.Mutex{}
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

	u := url.URL{Scheme: scheme, Host: host, Path: "/ws-api" + path}

	header := http.Header{}
	// header.Set("x-client", clientID)
	header.Set("Authorization", configuration.Connection.Token)

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	return conn, err
}

func WriteMessage(conn *websocket.Conn, messageId string, messageType entitie.MessageType, data any, errMsg string) (err error) {
	mu.Lock()
	defer mu.Unlock()
	var d []byte
	switch data.(type) {
	case []byte:
		d = data.([]byte)
	default:
		d, err = json.Marshal(data)
		if err != nil {
			return err
		}
	}

	body, err := json.Marshal(entitie.Message{
		MessageId:   messageId,
		MessageType: messageType,
		Data:        d,
		Error:       errMsg,
	})
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, body)
	if err != nil {
		return err
	}

	return nil
}
