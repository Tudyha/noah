package router

import (
	"noah/internal/server/controller"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func InitAdminRouter(api *gin.RouterGroup, i do.Injector) {
	adminGroup := api.Group("")
	adminController := do.MustInvoke[controller.AdminController](i)

	adminGroup.GET("/dashboard", adminController.Dashboard)
}
