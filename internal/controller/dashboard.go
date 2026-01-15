package controller

import (
	"noah/internal/service"
	"noah/pkg/response"

	"noah/pkg/sys"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type DashboardController struct {
	agentService service.AgentService
}

func NewDashboardController() *DashboardController {
	return &DashboardController{
		agentService: service.GetAgentService(),
	}
}

func (c *DashboardController) GetDashboard(ctx *gin.Context) {
	var response response.DashboardResponse
	basicInfo, _ := sys.GetBasicInfo()
	if basicInfo != nil {
		copier.Copy(&response.SysInfo, basicInfo)
	}
	hostInfo, _ := sys.GetHostInfo()
	if hostInfo != nil {
		copier.Copy(&response.SysInfo, hostInfo)
	}

	online, offline, _ := c.agentService.CountByAppID(ctx.Request.Context(), GetAppID(ctx))
	response.AgentStats.Online = online
	response.AgentStats.Offline = offline

	Success(ctx, response)
}
