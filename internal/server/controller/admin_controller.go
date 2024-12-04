package controller

import (
	"noah/internal/server/service"
	"noah/pkg/errcode"
	"noah/pkg/response"
	"noah/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/samber/do/v2"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type AdminController struct {
	clientService service.IClientService
}

func NewAdminController(i do.Injector) (AdminController, error) {
	return AdminController{
		clientService: do.MustInvoke[service.IClientService](i),
	}, nil
}

func (c AdminController) Dashboard(ctx *gin.Context) {
	var dashboardRes response.GetDashboardRes

	// 获取主机信息
	info, err := host.Info()
	if err != nil {
		Fail(ctx, errcode.ErrInternalError)
		return
	}
	copier.Copy(&dashboardRes, info)

	// 获取内存信息
	memStats, err := mem.VirtualMemory()
	if err != nil {
		Fail(ctx, errcode.ErrInternalError)
		return
	}
	dashboardRes.MemoryTotal = utils.CoverToGb(memStats.Total)
	dashboardRes.MemoryUsed = utils.CoverToGb(memStats.Used)
	dashboardRes.MemoryFree = utils.CoverToGb(memStats.Free)
	dashboardRes.MemoryAvailable = utils.CoverToGb(memStats.Available)
	dashboardRes.MemoryUsedPercent = utils.RoundToTwoDecimals(memStats.UsedPercent)
	dashboardRes.MemoryRemain = utils.RoundToTwoDecimals(dashboardRes.MemoryTotal - dashboardRes.MemoryUsed)

	// 获取磁盘信息
	usage, err := disk.Usage(".")
	if err != nil {
		Fail(ctx, errcode.ErrInternalError)
		return
	}
	dashboardRes.DiskTotal = utils.CoverToGb(usage.Total)
	dashboardRes.DiskUsed = utils.CoverToGb(usage.Used)
	dashboardRes.DiskFree = utils.CoverToGb(usage.Free)

	// 获取客户端信息
	online, offline := c.clientService.Count()
	dashboardRes.ClientOnlineCount = online
	dashboardRes.ClientOfflineCount = offline

	Success(ctx, dashboardRes)
}
