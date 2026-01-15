package v1

import (
	"noah/internal/controller"

	"github.com/gin-gonic/gin"
)

func RegisterDashboardRoutes(router *gin.RouterGroup) {
	dashboardController := controller.NewDashboardController()
	router.GET("/", dashboardController.GetDashboard)
}
