package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"noah/internal/server/controller"
	"noah/internal/server/middleware"
	"noah/internal/server/utils"
	"os"
	"path/filepath"
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

	clientController := handlers.GetClientController()
	ptyController := handlers.GetPtyController()
	fileController := handlers.GetFileController()

	{
		//免登录接口
		// 下载文件
		router.GET("/file/download/:filename", authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
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
		router.DELETE("/file/:filename", authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
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

	adminGroup := router.Group("", authMiddleware.MiddlewareFunc())
	{
		//需要登录接口
		clientGroup := adminGroup.Group("client")
		clientGroup.POST("/", clientController.CreateClient)
		clientGroup.GET("", clientController.GetClient)
		clientGroup.DELETE("/:id", clientController.DeleteClient)
		clientGroup.POST("/:id/health", handlers.Health)
		clientGroup.GET("/:id/ws", clientController.NewWsClient)
		clientGroup.POST("/cmd", clientController.SendCommandHandler)
		clientGroup.POST("/generate", clientController.Generate)
		clientGroup.POST("/:id/update", clientController.Update)
		clientGroup.GET("/:id/systemInfo", clientController.GetClientInfo)

		fileGroup := adminGroup.Group("/client/:id/file")
		fileGroup.GET("", fileController.GetFileList)
		fileGroup.GET("/content", fileController.GetFileContent)
		fileGroup.POST("/rename", fileController.RenameFile)
		fileGroup.DELETE("", fileController.DeleteFile)
		fileGroup.PUT("/content", fileController.UpdateFileContent)
		fileGroup.POST("", fileController.UploadFile)
		fileGroup.POST("/dir", fileController.NewDir)

		ptyGroup := adminGroup.Group("pty")
		ptyGroup.GET("/ws/:id", ptyController.NewPtyChannel)
		ptyGroup.GET("/client/ws/:channelId", ptyController.NewPtyClient)

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
	}

	router.Use(gin.Recovery())
}
