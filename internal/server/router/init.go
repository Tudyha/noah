package router

import (
	"noah/internal/server/controller"
	"noah/internal/server/middleware/auth"

	"github.com/gin-contrib/cors"
	"github.com/samber/do/v2"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine, i do.Injector) {
	// ----- CORS -----
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Token"}
	router.Use(cors.New(config))

	userController := do.MustInvoke[controller.UserController](i)
	authMiddleware := do.MustInvoke[*auth.AuthMiddleware](i)
	clientController := do.MustInvoke[controller.ClientController](i)

	//前端静态文件
	router.Static("/static", "web/dist/static") // 假设前端的静态资源在 /dist/assets 下
	router.GET("/", func(c *gin.Context) {
		c.File("web/dist/index.html")
	})
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.File("web/dist/favicon.ico")
	})

	api := router.Group("/api")
	api.POST("/login", userController.Login)

	api.POST("/client/connect", clientController.Connect)

	authApi := api.Group("", authMiddleware.AuthMiddlewareFunc())
	//临时授权接口
	tempAuthApi := api.Group("", authMiddleware.TempAuthMiddlewareFunc())
	tempAuthApi.GET("/client/build", clientController.GenerateClient)

	InitUserRouter(authApi, i)
	InitClientRouter(authApi, i)
	InitAdminRouter(authApi, i)
	InitTunnelRouter(authApi, i)
	InitCommonRouter(tempAuthApi, i)
}
