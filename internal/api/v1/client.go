package v1

import (
	"noah/internal/controller"

	"github.com/gin-gonic/gin"
)

func RegisterClientRoutes(router *gin.RouterGroup, clientController *controller.ClientController) {
	router.GET("/page", clientController.GetClientPage)
	router.GET("/bind", clientController.GetClientBind)
}
