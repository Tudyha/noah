package environment

import (
	"fmt"
	"strings"
)

func Load(serverAddress, httpPort, token string) *Configuration {
	return &Configuration{
		Connection: Connection{
			Token:           fmt.Sprint("Bearer ", token),
			ContextDeadline: 5,
		},
		Server: Server{
			Address:  serverAddress,
			HttpPort: httpPort,
			Url:      newServerUrl(serverAddress, httpPort),
		},
	}
}

func newServerUrl(serverAddress, serverPort string) string {
	if len(strings.TrimSpace(serverPort)) == 0 {
		return fmt.Sprintf("%s/api", strings.TrimRight(serverAddress, "/"))
	}
	return fmt.Sprintf("http://%s:%s/api", serverAddress, serverPort)
}
