package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"noah/internal/server/enum"
	"noah/internal/server/service"
	"os/exec"
	"strings"
	"sync"

	"noah/internal/server/dto"
	"noah/internal/server/utils"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type clientService struct {
	mu      *sync.Mutex
	clients map[uint]*websocket.Conn //客户端命令执行websocket
}

const (
	clientBaseDir  = "client/"
	buildBaseDir   = "build/"
	configFileName = "config.json"
	mainFileName   = "main.go"
	buildStr       = `GO_ENABLED=1 CGO_ENABLED=0 GOOS=%s GOARCH=amd64 go build -ldflags '%s -s -w -X main.Version=%s -extldflags "-static"' -o ../../temp/%s main.go`
)

func init() {
	service.RegisterClientService(&clientService{
		mu:      &sync.Mutex{},
		clients: make(map[uint]*websocket.Conn),
	})
}

var (
	ErrClientConnectionNotFound = errors.New("no active client connection found")
	ErrInvalidServerAddress     = errors.New("the server address provided is invalid")
	ErrInvalidServerPort        = errors.New("the server port provided is invalid")
)

func (c *clientService) AddConnection(id uint, connection *websocket.Conn) error {
	c.mu.Lock()
	c.clients[id] = connection
	c.mu.Unlock()
	return nil
}

func (c *clientService) getConnection(clientID uint) (*websocket.Conn, bool) {
	c.mu.Lock()
	conn, found := c.clients[clientID]
	c.mu.Unlock()
	return conn, found
}

// func (c *clientService) removeConnection(clientID uint) error {
// 	c.mu.Lock()
// 	delete(c.clients, clientID)
// 	c.mu.Unlock()
// 	return nil
// }

// SendCommand 执行命令
func (c *clientService) SendCommand(id uint, commandStr string, parameter string) (string, error) {
	client, found := c.getConnection(id)
	if !found {
		return ErrClientConnectionNotFound.Error(), nil
	}

	command := &dto.Command{
		Command:   commandStr,
		Parameter: parameter,
	}

	req, err := json.Marshal(command)
	if err != nil {
		return "", err
	}

	err = client.WriteMessage(websocket.BinaryMessage, req)
	if err != nil {
		return ErrClientConnectionNotFound.Error(), nil
	}

	_, readMessage, err := client.ReadMessage()
	if err != nil {
		return ErrClientConnectionNotFound.Error(), nil
	}

	var response dto.RespondCommandRequestBody
	if err := json.Unmarshal(readMessage, &response); err != nil {
		return "", err
	}

	command.Response = response.Response
	command.HasError = response.HasError

	command, err = handleResponse(command)
	if err != nil {
		return "", err
	}

	res := utils.ByteToString(command.Response)
	if command.HasError {
		return "", fmt.Errorf(res)
	}
	if len(strings.TrimSpace(res)) == 0 {
		return `No content.`, nil
	}
	return res, nil
}

func handleResponse(payload *dto.Command) (*dto.Command, error) {
	// const screenshotCmd = "screenshot"
	switch payload.Command {
	// case screenshotCmd:
	// 	filepath, err := image.WritePNG(payload.Response)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	payload.Response = utils.StringToByte(filepath)
	// 	break
	default:
		return payload, nil
	}
}

type ClientConfig struct {
	ServerAddress string `json:"server_address"`
	ServerPort    string `json:"server_port"`
}

func (c *clientService) Generate(serverAddr string, port string, osType int8, filename string) (string, error) {
	buildPath, err := c.PrepareBuildSession(serverAddr, port, osType)
	if err != nil {
		return "", err
	}
	defer utils.RemoveDir(buildPath)
	filename = buildFilename(enum.OSType(osType), filename)
	buildCmd := fmt.Sprintf(buildStr, getOSBuildParam(enum.OSType(osType)), getRunHiddenBuildParam(false), "dev", filename)

	cmd := exec.Command("sh", "-c", buildCmd)
	cmd.Dir = buildPath

	outputErr, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w:%s", err, outputErr)
	}
	return filename, nil
}

func (c clientService) BuildClientConfiguration(serverAddr string, port string) (clientConfig *ClientConfig, err error) {
	return &ClientConfig{
		ServerAddress: serverAddr,
		ServerPort:    port,
	}, err
}

func (c clientService) WriteClientConfigurationFile(clientConfig *ClientConfig, buildPath string) error {
	configurationJson, err := json.Marshal(clientConfig)
	if err != nil {
		return err
	}

	return utils.WriteFile(buildPath+configFileName, configurationJson)
}

func (c *clientService) PrepareBuildSession(serverAddr string, port string, osType int8) (string, error) {
	sessionID := uuid.New().String()
	buildPath := fmt.Sprint(buildBaseDir, sessionID, "/")

	err := utils.CopyDir(clientBaseDir, buildPath, configFileName)
	if err != nil {
		return "", err
	}

	clientConfiguration, err := c.BuildClientConfiguration(serverAddr, port)
	if err != nil {
		return "", err
	}

	err = c.WriteClientConfigurationFile(clientConfiguration, buildPath)
	if err != nil {
		return "", err
	}

	return buildPath, nil
}

func getOSBuildParam(osType enum.OSType) string {
	const (
		windowsKey = "windows"
		linuxKey   = "linux"
		macKey     = "darwin"
		unknownKey = "unknown"
	)
	switch osType {
	case enum.Windows:
		return windowsKey
	case enum.Linux:
		return linuxKey
	case enum.Darwin:
		return macKey
	default:
		return unknownKey
	}
}

func getRunHiddenBuildParam(hidden bool) string {
	if hidden {
		return "-H=windowsgui"
	}
	return ""
}

func buildFilename(os enum.OSType, filename string) string {
	if len(strings.TrimSpace(filename)) <= 0 {
		filename = uuid.New().String()
	}
	switch os {
	case enum.Windows:
		return fmt.Sprint(filename, ".exe")
	default:
		return filename
	}
}
