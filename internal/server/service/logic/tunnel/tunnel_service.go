package tunnel

import (
	"fmt"
	"net"
	"noah/internal/server/dao"
	"noah/internal/server/environment"
	"noah/internal/server/gateway"
	"noah/internal/server/middleware/log"
	"noah/internal/server/model"
	"noah/pkg/request"
	"noah/pkg/response"

	"github.com/jinzhu/copier"
	"github.com/samber/do/v2"
)

type tunnelService struct {
	tunnelDao dao.TunnelDao
	tunnels   map[uint]*tunnel
	gateway   *gateway.Gateway
	env       *environment.Environment
}

func NewTunnelService(i do.Injector) (tunnelService, error) {
	s := tunnelService{
		tunnelDao: do.MustInvoke[dao.TunnelDao](i),
		tunnels:   make(map[uint]*tunnel),
		gateway:   do.MustInvoke[*gateway.Gateway](i),
		env:       do.MustInvoke[*environment.Environment](i),
	}

	s.recoverTunnel()

	return s, nil
}

// NewTunnel 新建tunnel
func (c tunnelService) NewTunnel(id uint, tunnelReq request.CreateTunnelReq) error {
	tunnelType := tunnelReq.TunnelType
	serverPort := tunnelReq.ServerPort
	targetAddr := tunnelReq.TargetAddr
	cipher := tunnelReq.Cipher
	password := tunnelReq.Password

	if tunnelType == 2 {
		targetAddr = fmt.Sprintf("ss://%s:%s@%s:%d", cipher, password, c.env.Server.Host, serverPort)
	}

	tunnelId, err := c.tunnelDao.Save(model.Tunnel{
		TunnelType: tunnelType,
		ClientId:   id,
		ServerPort: serverPort,
		TargetAddr: targetAddr,
		Cipher:     cipher,
		Password:   password,
	})
	if err != nil {
		log.Error("Save tunnel error", map[string]interface{}{"clientId": id, "error": err})
		return err
	}
	if err = c.startTunnel(tunnelId); err != nil {
		return err
	}

	return nil
}

func (c tunnelService) startTunnel(tunnelId uint) (err error) {
	mt, err := c.tunnelDao.GetById(tunnelId)
	if err != nil {
		return err
	}
	t, err := newTunnel(mt.ID, mt.TunnelType, mt.ServerPort, mt.ClientId, mt.TargetAddr, mt.Cipher, mt.Password, &c)
	if err != nil {
		return err
	}
	if err = t.start(); err != nil {
		c.tunnelDao.UpdateStatus(tunnelId, 2, err.Error())
		return err
	} else {
		c.tunnelDao.UpdateStatus(tunnelId, 1, "")
	}
	c.tunnels[tunnelId] = t
	return err
}

func (c tunnelService) removeTunnel(id uint) {
	delete(c.tunnels, id)
}

// GetTunnelList 获取tunnel列表
func (c tunnelService) GetTunnelList(clientId uint) (res []response.GetTunnelListRes, err error) {
	list, err := c.tunnelDao.List(clientId)
	if err != nil {
		return nil, err
	}
	copier.Copy(&res, list)
	return res, nil
}

// DeleteTunnel 删除tunnel
func (c tunnelService) DeleteTunnel(id uint) error {
	_, err := c.tunnelDao.GetById(id)
	if err != nil {
		return err
	}
	err = c.tunnelDao.Delete(id)
	if err != nil {
		return err
	}

	// 关闭端口监听
	if t, ok := c.tunnels[id]; ok {
		t.stop()
		c.removeTunnel(id)
	}
	return nil
}

// recoverTunnel 恢复tunnel
func (c tunnelService) recoverTunnel() {
	list, err := c.tunnelDao.List(0)
	if err != nil {
		return
	}

	for _, tunnel := range list {
		c.startTunnel(tunnel.ID)
	}
}

func (c tunnelService) newClientConn(clientId uint, network, targetAddr string) (net.Conn, error) {
	return c.gateway.NewClientConn(uint32(clientId), network, targetAddr)
}
