package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"noah/internal/server/controller"
	"noah/internal/server/middleware"
)

type Router struct {
	G        *gin.Engine
	handlers *controller.Controller
}

func NewRouter(g *gin.Engine) *Router {
	return &Router{
		G:        g,
		handlers: controller.NewController(),
	}
}

func (r *Router) LoadRoutes() {
	router := r.G

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
	router.POST("/login", authMiddleware.LoginHandler)

	auth := router.Group("/auth", authMiddleware.MiddlewareFunc())
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)

	// -----  jwt  -----

	handlers := r.handlers

	deviceController := handlers.GetDeviceController()
	clientController := handlers.GetClientController()
	//shellController := handlers.GetShellController()
	ptyController := handlers.GetPtyController()
	fileController := handlers.GetFileController()

	{
		//免登录接口
		clientGroup := router.Group("client")
		clientGroup.GET("/health/:id", handlers.Health)
		clientGroup.POST("/device", deviceController.CreateDevice)
		clientGroup.GET("/ws/:id", clientController.NewClient)

		ptyGroup := router.Group("pty")
		ptyGroup.GET("/ws/:id", ptyController.WebSocket)
		ptyGroup.GET("/client/ws/:channelId", clientController.NewPtyClient)

		router.GET("/download/:filename", clientController.Download)

	}

	adminGroup := router.Group("", authMiddleware.MiddlewareFunc())
	{
		//需要登录接口
		deviceGroup := adminGroup.Group("device")
		deviceGroup.GET("", deviceController.GetDevice)
		deviceGroup.DELETE("", deviceController.DeleteDevice)

		//shellGroup := router.Group("shell")
		//shellGroup.GET("/ws/:id", shellController.WebSocket)

		adminClientGroup := adminGroup.Group("client")
		adminClientGroup.POST("/cmd", clientController.SendCommandHandler)
		adminClientGroup.POST("/generate", clientController.Generate)

		userGroup := adminGroup.Group("user")
		//userGroup.POST("login", func(ctx *gin.Context) {
		//	ctx.JSON(http.StatusOK, gin.H{
		//		"code": 0,
		//		"data": "admin-token",
		//	})
		//})
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

		fileGroup := adminGroup.Group("/client/:id/file")
		fileGroup.GET("", fileController.GetFileList)
		fileGroup.GET("/content", fileController.GetFileContent)
		fileGroup.POST("/rename", fileController.RenameFile)
		fileGroup.DELETE("", fileController.DeleteFile)
		fileGroup.PUT("/content", fileController.UpdateFileContent)
		fileGroup.POST("", fileController.UploadFile)
	}

	router.Use(gin.Recovery())
}
