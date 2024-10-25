package routes

import (
	"log"
	"net/http"
	"noah/internal/server/controller"
	"noah/internal/server/middleware"
	"noah/internal/server/utils"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/samber/do/v2"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Gin      *gin.Engine
	handlers *controller.Controller
	i        do.Injector
}

func NewRouter(g *gin.Engine, i do.Injector) *Router {
	return &Router{
		Gin:      g,
		handlers: controller.NewController(i),
		i:        i,
	}
}

func (r *Router) LoadRoutes() {
	router := r.Gin

	// ----- CORS -----
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Token"}
	router.Use(cors.New(config))

	authMiddleware := middleware.NewAuthMiddleware(r.i)
	do.ProvideValue(r.i, authMiddleware)

	handlers := r.handlers
	clientController := handlers.GetClientController()
	channelController := handlers.GetChannelController()
	fileController := handlers.GetFileController()
	adminController := handlers.GetAdminController()
	userController := handlers.GetUserController()

	{
		//前端静态文件
		router.Static("/static", "web/dist/static") // 假设前端的静态资源在 /dist/assets 下
		router.GET("/", func(c *gin.Context) {
			c.File("web/dist/index.html")
		})
		router.GET("/favicon.ico", func(c *gin.Context) {
			c.File("web/dist/favicon.ico")
		})
	}

	api := router.Group("/api")
	{
		//免登录接口
		api.POST("/login", userController.Login)
		api.POST("/refresh_token", userController.RefreshToken)
	}

	// 客户端接口
	{
		clientGroup := api.Group("client", middleware.ClientAuthMiddlewareFunc())
		clientGroup.POST("/", clientController.CreateClient)
		clientGroup.POST("/:id/health", clientController.Health)
	}

	// 管理员接口
	{
		adminGroup := api.Group("", middleware.AdminAuthMiddlewareFunc())
		clientGroup := adminGroup.Group("client")
		clientGroup.GET("/:id", clientController.GetClient)
		clientGroup.GET("/page", clientController.GetClientPage)
		clientGroup.DELETE("/:id", clientController.DeleteClient)
		clientGroup.POST("/cmd", clientController.SendCommandHandler)
		clientGroup.POST("/generate", clientController.Generate)
		clientGroup.POST("/:id/update", clientController.Update)
		clientGroup.GET("/:id/systemInfo", clientController.GetClientInfo)
		clientGroup.GET("/:id/process", clientController.GetClientProcessList)
		clientGroup.DELETE("/:id/process/:pid", clientController.KillClientProcess)
		clientGroup.GET("/:id/network", clientController.GetClientNetworkList)
		clientGroup.GET("/:id/docker/container", clientController.GetClientDockerContainerList)

		fileGroup := adminGroup.Group("/client/:id/file")
		fileGroup.GET("", fileController.GetFileList)
		fileGroup.GET("/content", fileController.GetFileContent)
		fileGroup.POST("/rename", fileController.RenameFile)
		fileGroup.DELETE("", fileController.DeleteFile)
		fileGroup.PUT("/content", fileController.UpdateFileContent)
		fileGroup.POST("", fileController.UploadFile)
		fileGroup.POST("/dir", fileController.NewDir)

		userGroup := adminGroup.Group("user")
		userGroup.GET("info", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"code": 0,
				"data": gin.H{
					"roles":        []string{"admin"},
					"introduction": "I am a super administrator",
					"avatar":       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
					"name":         "Super Admin",
				},
			})
		})

		chGroup := adminGroup.Group("/client/:id/channel")
		chGroup.POST("", channelController.NewChannel)
		chGroup.GET("", channelController.GetChannelList)
		chGroup.DELETE("/:channelId", channelController.DeleteChannel)

		// dashboard
		adminGroup.GET("/dashboard", adminController.Dashboard)
		adminGroup.GET("/generateClientToken", adminController.GenerateClientToken)
	}
	{
		// 通用接口
		commonGroup := api.Group("", middleware.AuthMiddlewareFunc())
		// 下载文件
		commonGroup.GET("/file/download/:filename", func(c *gin.Context) {
			filename := c.Param("filename")
			sanitizedFilename := filepath.Base(filename)
			filePath := "temp/" + sanitizedFilename

			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				c.AbortWithStatus(http.StatusNotFound)
				log.Printf("File not found: %v", err)
				return
			} else if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				log.Printf("Error while accessing file: %v", err)
				return
			}

			c.File(filePath)
		})

		// 删除文件
		commonGroup.DELETE("/file/:filename", func(c *gin.Context) {
			filename := c.Param("filename")
			sanitizedFilename := filepath.Base(filename)
			filePath := "temp/" + sanitizedFilename

			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				c.AbortWithStatus(http.StatusNotFound)
				log.Printf("File not found: %v", err)
				return
			} else if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				log.Printf("Error while accessing file: %v", err)
				return
			}

			if err := utils.RemoveFile(filePath); err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				log.Printf("Error while removing file: %v", err)
				return
			}

			c.Status(http.StatusOK)
		})
	}

	// websocket接口
	wsApi := router.Group("/ws-api", middleware.AuthMiddlewareFunc())
	{
		wsApi.GET("/client/:id/ws", clientController.NewWsClient)
		ptyGroup := wsApi.Group("pty")
		ptyGroup.GET("/ws/:id", channelController.NewPtyChannel)
	}

	router.Use(gin.Recovery())
}
