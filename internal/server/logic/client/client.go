package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"noah/internal/server/dao"
	"noah/internal/server/enum"
	"noah/internal/server/vo"
	"os/exec"
	"strings"
	"sync"
	"time"

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

func NewClientService() *clientService {
	return &clientService{
		mu:      &sync.Mutex{},
		clients: make(map[uint]*websocket.Conn),
	}
}

var (
	ErrClientConnectionNotFound = errors.New("no active client connection found")
	ErrInvalidServerAddress     = errors.New("the server address provided is invalid")
	ErrInvalidServerPort        = errors.New("the server port provided is invalid")
)

// AddConnection 新增ws连接
func (c clientService) AddConnection(id uint, connection *websocket.Conn) error {
	// 验证连接是否有效
	if connection == nil || connection.RemoteAddr() == nil {
		return errors.New("invalid connection")
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	// 更新或添加新连接
	c.clients[id] = connection
	return nil
}

// getConnection 获取连接
func (c clientService) getConnection(clientID uint) (*websocket.Conn, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	conn, found := c.clients[clientID]
	return conn, found
}

// removeConnection 删除连接
func (c clientService) removeConnection(clientID uint) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if conn, found := c.clients[clientID]; !found {
	} else {
		err := conn.Close()
		if err != nil {
		}
	}
	delete(c.clients, clientID)
	return nil
}

// SendCommand 执行命令
func (c clientService) SendCommand(id uint, commandStr string, parameter string) (string, error) {
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
	//if len(strings.TrimSpace(res)) == 0 {
	//	return `No content.`, nil
	//}
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
	Token         string `json:"token"`
}

func (c clientService) Generate(serverAddr string, port string, osType int8, token string, filename string) (string, error) {
	buildPath, err := c.PrepareBuildSession(serverAddr, port, token)
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

func (c clientService) BuildClientConfiguration(serverAddr string, port string, token string) (clientConfig *ClientConfig, err error) {
	return &ClientConfig{
		ServerAddress: serverAddr,
		ServerPort:    port,
		Token:         token,
	}, err
}

func (c clientService) WriteClientConfigurationFile(clientConfig *ClientConfig, buildPath string) error {
	configurationJson, err := json.Marshal(clientConfig)
	if err != nil {
		return err
	}

	return utils.WriteFile(buildPath+configFileName, configurationJson)
}

func (c clientService) PrepareBuildSession(serverAddr string, port string, token string) (string, error) {
	sessionID := uuid.New().String()
	buildPath := fmt.Sprint(buildBaseDir, sessionID, "/")

	err := utils.CopyDir(clientBaseDir, buildPath, configFileName)
	if err != nil {
		return "", err
	}

	clientConfiguration, err := c.BuildClientConfiguration(serverAddr, port, token)
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

// Exit 关闭连接
func (c clientService) Exit(id uint) error {
	if err := c.removeConnection(id); err != nil {
		return err
	}
	return nil
}

func (c clientService) Save(body dto.ClientPostDto) (id uint, err error) {
	var Client dao.Client
	copier.Copy(&Client, &body)
	Client.OsType = enum.DetectOS(Client.OSName)

	old := dao.GetClientDao().GetByMacAddress(Client.MacAddress)
	if old.ID != 0 {
		// 已存在，更新数据即可
		Client.ID = old.ID
		err = dao.GetClientDao().Update(Client)
		if err != nil {
			return 0, err
		}
		return old.ID, nil
	}

	id, err = dao.GetClientDao().Save(Client)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c clientService) UpdateStatus(id uint, status int8) {
	dao.GetClientDao().UpdateStatus(id, status)
}

func (c clientService) GetClient(query vo.ClientListQueryReq) (total int64, ClientDtos []dto.ClientDto) {
	total, Clients := dao.GetClientDao().Page(query)

	if total == 0 {
		return 0, make([]dto.ClientDto, 0)
	}

	copier.Copy(&ClientDtos, &Clients)
	return total, ClientDtos
}

func (c clientService) ScheduleUpdateStatus() error {
	dao.GetClientDao().ScheduleUpdateStatus()
	return nil
}

func (c clientService) Delete(id uint) error {
	return dao.GetClientDao().Delete(id)
}

func (c clientService) SaveSystemInfo(id uint, systemInfo dto.SystemInfoReq) error {
	var clientInfo dao.ClientInfo
	copier.Copy(&clientInfo, systemInfo)
	clientInfo.ClientID = id
	return dao.GetClientInfoDao().Create(clientInfo)
}

func (c clientService) GetSystemInfo(id uint, start time.Time, end time.Time) ([]dto.SystemInfoRes, error) {
	clientInfoList := dao.GetClientInfoDao().GetByClientId(id, start, end)
	if len(clientInfoList) == 0 {
		return make([]dto.SystemInfoRes, 0), nil
	}
	var result []dto.SystemInfoRes
	copier.Copy(&result, clientInfoList)
	return result, nil
}

func (c clientService) CleanSystemInfo() error {
	dao.GetClientInfoDao().Clean()
	return nil
}
