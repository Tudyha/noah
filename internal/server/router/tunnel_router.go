package router

import (
	"noah/internal/server/controller"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func InitTunnelRouter(api *gin.RouterGroup, i do.Injector) {
	g := api.Group("/client/:id/tunnel")
	tunnelController := do.MustInvoke[controller.TunnelController](i)

	g.POST("", tunnelController.NewTunnel)
	g.DELETE("/:tunnelId", tunnelController.DeleteTunnel)
	g.GET("", tunnelController.GetTunnelList)
}
