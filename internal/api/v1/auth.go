package v1

import (
	"noah/internal/controller"

	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes 注册认证相关路由
func RegisterAuthRoutes(router *gin.RouterGroup, authController *controller.AuthController) {
	router.POST("/login", authController.Login)
}
