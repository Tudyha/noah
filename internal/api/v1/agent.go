package v1

import (
	"noah/internal/controller"

	"github.com/gin-gonic/gin"
)

func RegisterAgentRoutes(router *gin.RouterGroup, agentController *controller.AgentController) {
	router.GET("/page", agentController.GetAgentPage)
	router.GET("/bind", agentController.GetAgentBind)
	router.DELETE("/:agent_id", agentController.DeleteAgent)
	router.GET("/:agent_id/metric", agentController.GetAgentMetric)
	router.GET("/:agent_id/pty", agentController.OpenPty)
	router.GET("/:agent_id", agentController.GetAgent)
	router.POST("/v2ray/sub", agentController.GenerateV2raySubscribeLink)
}
