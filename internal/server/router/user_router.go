package router

import (
	"noah/internal/server/controller"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func InitUserRouter(api *gin.RouterGroup, i do.Injector) {
	userGroup := api.Group("user")
	userController := do.MustInvoke[controller.UserController](i)

	userGroup.GET("/info", userController.GetUser)
	userGroup.GET("/page", userController.GetUserPage)
	userGroup.PATCH("/:userId/password", userController.UpdatePassword)
}
