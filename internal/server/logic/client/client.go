package client

import (
	"encoding/json"
	"fmt"
	"noah/internal/server/dao"
	"noah/internal/server/enum"
	"noah/internal/server/request"
	"noah/internal/server/response"
	"os/exec"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"github.com/samber/do/v2"

	"noah/internal/server/utils"

	"github.com/google/uuid"
)

type Service struct {
	clientDao     *dao.ClientDao
	clientInfoDao *dao.ClientInfoDao
}

const (
	clientBaseDir  = "client/"
	buildBaseDir   = "build/"
	configFileName = "config.json"
	pkgDir         = "pkg/"
	buildStr       = `CGO_ENABLED=0 GOOS=%s GOARCH=amd64 go build -ldflags '%s -s -w -X main.Version=%s -extldflags "-static"' -o ../../temp/%s main.go`
)

func NewClientService(i do.Injector) *Service {
	return &Service{
		clientDao:     do.MustInvoke[*dao.ClientDao](i),
		clientInfoDao: do.MustInvoke[*dao.ClientInfoDao](i),
	}
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
	//defer utils.RemoveDir(buildPath)
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

	utils.CopyDir(pkgDir, buildBaseDir+"pkg/", configFileName)

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

	old := c.clientDao.GetByMacAddress(client.MacAddress)
	if old.ID != 0 {
		// 已存在，更新数据即可
		client.ID = old.ID
		err = c.clientDao.Update(client)
		if err != nil {
			return 0, err
		}
		return old.ID, nil
	}

	id, err = c.clientDao.Save(client)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c Service) UpdateStatus(id uint, status int8) {
	c.clientDao.UpdateStatus(id, status)
}

func (c Service) GetClientPage(query request.ListClientQueryReq) (int64, []response.ListClientRes) {
	total, clients := c.clientDao.Page(query)

	if total == 0 {
		return 0, make([]response.ListClientRes, 0)
	}

	var res []response.ListClientRes
	copier.Copy(&res, &clients)
	return total, res
}

func (c Service) GetClient(id uint) (response.GetClientRes, error) {
	client, err := c.clientDao.GetById(id)
	if err != nil {
		return response.GetClientRes{}, err
	}
	var res response.GetClientRes
	copier.Copy(&res, client)

	res.MemoryTotal = fmt.Sprintf("%.0f GB", utils.CoverToGb(client.MemoryTotal))
	res.DiskTotal = fmt.Sprintf("%.0f GB", utils.CoverToGb(client.DiskTotal))
	return res, nil
}

func (c Service) ScheduleUpdateStatus() error {
	c.clientDao.ScheduleUpdateStatus()
	return nil
}

func (c Service) Delete(id uint) error {
	return c.clientDao.Delete(id)
}

func (c Service) SaveSystemInfo(id uint, systemInfo request.CreateSystemInfoReq) error {
	var clientInfo dao.ClientInfo
	copier.Copy(&clientInfo, systemInfo)
	clientInfo.ClientID = id
	return c.clientInfoDao.Create(clientInfo)
}

func (c Service) GetSystemInfo(id uint, start time.Time, end time.Time) ([]response.GetSystemInfoRes, error) {
	clientInfoList := c.clientInfoDao.GetByClientId(id, start, end)
	if len(clientInfoList) == 0 {
		return make([]response.GetSystemInfoRes, 0), nil
	}
	result := make([]response.GetSystemInfoRes, 0)
	for _, clientInfo := range clientInfoList {
		var res response.GetSystemInfoRes
		copier.Copy(&res, clientInfo)

		res.MemoryTotal = utils.CoverToGb(clientInfo.MemoryTotal)
		res.MemoryFree = utils.CoverToGb(clientInfo.MemoryFree)
		res.MemoryUsed = utils.CoverToGb(clientInfo.MemoryUsed)
		res.MemoryAvailable = utils.CoverToGb(clientInfo.MemoryAvailable)
		res.DiskTotal = utils.CoverToGb(clientInfo.DiskTotal)
		res.DiskFree = utils.CoverToGb(clientInfo.DiskFree)
		res.DiskUsed = utils.CoverToGb(clientInfo.DiskUsed)

		res.BandwidthIn = utils.CoverToKb(clientInfo.BandwidthIn)
		res.BandwidthOut = utils.CoverToKb(clientInfo.BandwidthOut)
		result = append(result, res)
	}

	return result, nil
}

func (c Service) CleanSystemInfo() error {
	c.clientInfoDao.Clean()
	return nil
}

func (c Service) Count() (online int64, offline int64) {
	return c.clientDao.Count()
}
