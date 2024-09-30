package routes

import (
	"log"
	"net/http"
	"noah/internal/server/controller"
	"noah/internal/server/gateway"
	"noah/internal/server/middleware"
	"noah/internal/server/utils"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Gin      *gin.Engine
	handlers *controller.Controller
	gateway  *gateway.Gateway
}

func NewRouter(g *gin.Engine, gate *gateway.Gateway) *Router {
	return &Router{
		Gin:      g,
		handlers: controller.NewController(gate),
		gateway:  gate,
	}
}

func (r *Router) LoadRoutes() {
	router := r.Gin

	// ----- CORS -----
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(config))
	// ----- CORS -----

	// -----  jwt  -----
	authMiddleware, err := middleware.RegisterJwtMiddleWare()
	if err != nil {
		panic(err)
	}
	router.Use(func(context *gin.Context) {
		errInit := authMiddleware.MiddlewareInit()
		if errInit != nil {
			log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
		}
	})
	// -----  jwt  -----

	handlers := r.handlers

	clientController := handlers.GetClientController()
	channelController := handlers.GetChannelController()
	fileController := handlers.GetFileController()
	adminController := handlers.GetAdminController()

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
		api.POST("/login", authMiddleware.LoginHandler)
		api.GET("/refresh_token", authMiddleware.RefreshHandler)
	}

	apiAuthGroup := api.Group("", authMiddleware.MiddlewareFunc())
	{
		clientGroup := apiAuthGroup.Group("client")
		clientGroup.POST("/", clientController.CreateClient)
		clientGroup.GET("/:id", clientController.GetClient)
		clientGroup.GET("/page", clientController.GetClientPage)
		clientGroup.DELETE("/:id", clientController.DeleteClient)
		clientGroup.POST("/:id/health", clientController.Health)
		clientGroup.POST("/cmd", clientController.SendCommandHandler)
		clientGroup.POST("/generate", clientController.Generate)
		clientGroup.POST("/:id/update", clientController.Update)
		clientGroup.GET("/:id/systemInfo", clientController.GetClientInfo)
		clientGroup.GET("/:id/process", clientController.GetClientProcessList)
		clientGroup.DELETE("/:id/process/:pid", clientController.KillClientProcess)
		clientGroup.GET("/:id/network", clientController.GetClientNetworkList)
		clientGroup.GET("/:id/docker/container", clientController.GetClientDockerContainerList)

		fileGroup := apiAuthGroup.Group("/client/:id/file")
		fileGroup.GET("", fileController.GetFileList)
		fileGroup.GET("/content", fileController.GetFileContent)
		fileGroup.POST("/rename", fileController.RenameFile)
		fileGroup.DELETE("", fileController.DeleteFile)
		fileGroup.PUT("/content", fileController.UpdateFileContent)
		fileGroup.POST("", fileController.UploadFile)
		fileGroup.POST("/dir", fileController.NewDir)

		userGroup := apiAuthGroup.Group("user")
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

		chGroup := apiAuthGroup.Group("/client/:id/channel")
		chGroup.POST("", channelController.NewChannel)
		chGroup.GET("", channelController.GetChannelList)
		chGroup.DELETE("/:channelId", channelController.DeleteChannel)

		// 下载文件
		apiAuthGroup.GET("/file/download/:filename", func(c *gin.Context) {
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
		apiAuthGroup.DELETE("/file/:filename", func(c *gin.Context) {
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

		//dashboard
		apiAuthGroup.GET("/dashboard", adminController.Dashboard)
	}

	// websocket接口
	wsApi := router.Group("/ws-api", authMiddleware.MiddlewareFunc())
	{
		wsApi.GET("/client/:id/ws", clientController.NewWsClient)
		ptyGroup := wsApi.Group("pty")
		ptyGroup.GET("/ws/:id", channelController.NewPtyChannel)
	}

	router.Use(gin.Recovery())
}
