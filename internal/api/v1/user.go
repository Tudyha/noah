package v1

import (
	"noah/internal/controller"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 注册用户相关路由
func RegisterUserRoutes(router *gin.RouterGroup, userController *controller.UserController) {
	router.GET("/", userController.GetUser)
}
