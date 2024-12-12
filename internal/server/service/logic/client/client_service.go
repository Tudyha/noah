package client

import (
	"encoding/json"
	"fmt"
	"noah/internal/server/dao"
	"noah/internal/server/gateway"
	"noah/internal/server/model"
	"noah/pkg/enum"
	"noah/pkg/request"
	"noah/pkg/response"
	"noah/pkg/utils"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/robfig/cron/v3"
	"github.com/samber/do/v2"
)

var (
	clientBasePath       = "client/"
	clientConfigFilename = "config.json"
	buildStr             = `CGO_ENABLED=0 GOOS=%s GOARCH=%s go build -ldflags ' -s -w -extldflags "-static"' -o %s main.go`
)

type clientService struct {
	clientDao     dao.ClientDao
	clientStatDao dao.ClientStatDao
	gateway       *gateway.Gateway
}

func NewClientService(i do.Injector) (clientService, error) {
	s := clientService{
		clientDao:     do.MustInvoke[dao.ClientDao](i),
		clientStatDao: do.MustInvoke[dao.ClientStatDao](i),
		gateway:       do.MustInvoke[*gateway.Gateway](i),
	}

	s.gateway.SetPongHandler(func(clientId uint32, data []byte) {
		var clientStat request.CreateClientStatReq
		err := json.Unmarshal(data, &clientStat)
		if err != nil {
			return
		}
		s.SaveClientStat(uint(clientId), clientStat)
		s.clientDao.UpdateStatus(uint(clientId), enum.DEVICE_ONLINE)
	})

	s.runScheduleTask()

	return s, nil
}

func (c clientService) Save(client model.Client) (id uint, err error) {

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

func (c clientService) GetClientPage(query request.ListClientQueryReq) (int64, []response.ListClientRes) {
	total, clients := c.clientDao.Page(query.Hostname, query.Status, query.Page, query.Size)

	if total == 0 {
		return 0, make([]response.ListClientRes, 0)
	}

	var res []response.ListClientRes
	copier.Copy(&res, &clients)
	return total, res
}

func (c clientService) GetClient(id uint) (response.GetClientRes, error) {
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

func (c clientService) ScheduleUpdateStatus() error {
	c.clientDao.ScheduleUpdateStatus()
	return nil
}

func (c clientService) Delete(id uint) error {
	return c.clientDao.Delete(id)
}

func (c clientService) SaveClientStat(id uint, systemInfo request.CreateClientStatReq) error {
	var clientStat model.ClientStat
	copier.Copy(&clientStat, systemInfo)
	clientStat.ClientID = id
	return c.clientStatDao.Create(clientStat)
}

func (c clientService) GetClientStat(id uint, start time.Time, end time.Time) ([]response.GetClientStatRes, error) {
	clientInfoList := c.clientStatDao.GetByClientId(id, start, end)
	if len(clientInfoList) == 0 {
		return make([]response.GetClientStatRes, 0), nil
	}
	result := make([]response.GetClientStatRes, 0)
	for _, clientInfo := range clientInfoList {
		var res response.GetClientStatRes
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

func (c clientService) Count() (online int64, offline int64) {
	return c.clientDao.Count()
}

func (c clientService) runScheduleTask() {
	cr := cron.New()
	_, err := cr.AddFunc("* * * * *", func() {
		c.clientDao.ScheduleUpdateStatus()
	})
	if err != nil {
		return
	}

	//每天0点执行一次
	_, err = cr.AddFunc("0 0 * * *", func() {
		c.clientStatDao.Clean()
	})

	if err != nil {
		return
	}

	cr.Start()
}

func (c clientService) BuildCllient(goos, goarch, host, port string, compress uint8) (file string, err error) {
	filename, err := prepareBuildSession(host, port, compress)
	defer func() {
		os.Rename(clientBasePath+"/"+clientConfigFilename+"."+"b", clientBasePath+clientConfigFilename)
	}()

	file = filename

	buildCmd := fmt.Sprintf(buildStr, goos, goarch, file)

	cmd := exec.Command("sh", "-c", buildCmd)
	cmd.Dir = clientBasePath

	_, err = cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return clientBasePath + file, nil
}

func prepareBuildSession(host, port string, compress uint8) (filename string, err error) {
	sessionID := uuid.New().String()

	// 备份配置文件, 防止被覆盖
	if err = os.Rename(clientBasePath+clientConfigFilename, clientBasePath+"/"+clientConfigFilename+"."+"b"); err != nil {
		return
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		return
	}
	clientConfig := map[string]interface{}{
		"server": map[string]interface{}{
			"host": host,
			"port": p,
		},
		"compress": compress,
	}
	d, err := json.Marshal(clientConfig)
	if err != nil {
		return
	}
	if err = os.WriteFile(clientBasePath+clientConfigFilename, d, 0644); err != nil {
		return
	}
	return sessionID, nil
}
