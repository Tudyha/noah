package client

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"noah/internal/server/dao"
	"noah/internal/server/enum"
	"noah/internal/server/request"
	"noah/internal/server/response"
	"os/exec"
	"strings"
	"time"

	"noah/internal/server/dto"
	"noah/internal/server/utils"

	"github.com/google/uuid"
)

type Service struct {
}

const (
	clientBaseDir  = "client/"
	buildBaseDir   = "build/"
	configFileName = "config.json"
	buildStr       = `CGO_ENABLED=0 GOOS=%s GOARCH=amd64 go build -ldflags '%s -s -w -X main.Version=%s -extldflags "-static"' -o ../../temp/%s main.go`
)

func NewClientService() *Service {
	return &Service{}
}

type ClientConfig struct {
	ServerAddress string `json:"server_address"`
	ServerPort    string `json:"server_port"`
	Token         string `json:"token"`
}

func (c Service) Generate(serverAddr string, port string, osType int8, token string, filename string) (string, error) {
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

func (c Service) BuildClientConfiguration(serverAddr string, port string, token string) (clientConfig *ClientConfig, err error) {
	return &ClientConfig{
		ServerAddress: serverAddr,
		ServerPort:    port,
		Token:         token,
	}, err
}

func (c Service) WriteClientConfigurationFile(clientConfig *ClientConfig, buildPath string) error {
	configurationJson, err := json.Marshal(clientConfig)
	if err != nil {
		return err
	}

	return utils.WriteFile(buildPath+configFileName, configurationJson)
}

func (c Service) PrepareBuildSession(serverAddr string, port string, token string) (string, error) {
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

func (c Service) Save(client dao.Client) (id uint, err error) {

	old := dao.GetClientDao().GetByMacAddress(client.MacAddress)
	if old.ID != 0 {
		// 已存在，更新数据即可
		client.ID = old.ID
		err = dao.GetClientDao().Update(client)
		if err != nil {
			return 0, err
		}
		return old.ID, nil
	}

	id, err = dao.GetClientDao().Save(client)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c Service) UpdateStatus(id uint, status int8) {
	dao.GetClientDao().UpdateStatus(id, status)
}

func (c Service) GetClientPage(query request.ListClientQueryReq) (int64, []response.ListClientRes) {
	total, clients := dao.GetClientDao().Page(query)

	if total == 0 {
		return 0, make([]response.ListClientRes, 0)
	}

	var res []response.ListClientRes
	copier.Copy(&res, &clients)
	return total, res
}

func (c Service) GetClient(id uint) (response.GetClientRes, error) {
	client, err := dao.GetClientDao().GetById(id)
	if err != nil {
		return response.GetClientRes{}, err
	}
	var res response.GetClientRes
	copier.Copy(&res, client)
	return res, nil
}

func (c Service) ScheduleUpdateStatus() error {
	dao.GetClientDao().ScheduleUpdateStatus()
	return nil
}

func (c Service) Delete(id uint) error {
	return dao.GetClientDao().Delete(id)
}

func (c Service) SaveSystemInfo(id uint, systemInfo dto.SystemInfoReq) error {
	var clientInfo dao.ClientInfo
	copier.Copy(&clientInfo, systemInfo)
	clientInfo.ClientID = id
	return dao.GetClientInfoDao().Create(clientInfo)
}

func (c Service) GetSystemInfo(id uint, start time.Time, end time.Time) ([]dto.SystemInfoRes, error) {
	clientInfoList := dao.GetClientInfoDao().GetByClientId(id, start, end)
	if len(clientInfoList) == 0 {
		return make([]dto.SystemInfoRes, 0), nil
	}
	var result []dto.SystemInfoRes
	copier.Copy(&result, clientInfoList)
	return result, nil
}

func (c Service) CleanSystemInfo() error {
	dao.GetClientInfoDao().Clean()
	return nil
}
