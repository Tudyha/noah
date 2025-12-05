package api

import (
	v1 "noah/internal/api/v1"
	"noah/internal/controller"
	"noah/pkg/response"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	_ "noah/docs" // 导入生成的 docs 包
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(router *gin.Engine) {
	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, response.Response{
			Code: 0,
			Msg:  "OK",
		})
	})

	// 添加 Swagger UI 路由
	// 访问路径通常是 /swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1版本
	v1Group := router.Group("/api/v1")
	{

		authController := controller.GetAuthController()
		userController := controller.GetUserController()
		// 认证相关路由
		authGroup := v1Group.Group("/auth")
		{
			v1.RegisterAuthRoutes(authGroup, authController)
		}

		// 用户相关路由
		userGroup := v1Group.Group("/user")
		{
			v1.RegisterUserRoutes(userGroup, userController)
		}
	}
}
