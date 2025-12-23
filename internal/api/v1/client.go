package v1

import (
	"noah/internal/controller"

	"github.com/gin-gonic/gin"
)

func RegisterClientRoutes(router *gin.RouterGroup, clientController *controller.ClientController) {
	router.GET("/page", clientController.GetClientPage)
	router.GET("/bind", clientController.GetClientBind)
	router.DELETE("/:client_id", clientController.DeleteClient)
	router.GET("/:client_id/stat", clientController.GetClientStat)
	router.GET("/:client_id/pty", clientController.OpenPty)
}
