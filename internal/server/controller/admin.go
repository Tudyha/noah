package controller

import (
	"net/http"
	"noah/internal/server/middleware"
	"noah/internal/server/response"
	"noah/internal/server/service"
	"noah/internal/server/utils"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type AdminController struct {
}

func NewAdminController() *AdminController {
	return &AdminController{}
}

func (c AdminController) Dashboard(ctx *gin.Context) {
	var dashboardRes response.GetDashboardRes

	// 获取主机信息
	info, err := host.Info()
	if err != nil {
		Fail(ctx, http.StatusInternalServerError, err.Error())
	}
	copier.Copy(&dashboardRes, info)

	// 获取内存信息
	memStats, err := mem.VirtualMemory()
	if err != nil {
		Fail(ctx, http.StatusInternalServerError, err.Error())
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
		Fail(ctx, http.StatusInternalServerError, err.Error())
	}
	dashboardRes.DiskTotal = utils.CoverToGb(usage.Total)
	dashboardRes.DiskUsed = utils.CoverToGb(usage.Used)
	dashboardRes.DiskFree = utils.CoverToGb(usage.Free)

	// 获取客户端信息
	online, offline := service.GetClientService().Count()
	dashboardRes.ClientOnlineCount = online
	dashboardRes.ClientOfflineCount = offline

	Success(ctx, dashboardRes)
}

func (c AdminController) GenerateClientToken(ctx *gin.Context) {
	Success(ctx, middleware.GenerateClientToken())
}
