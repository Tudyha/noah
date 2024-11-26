package router

import (
	"noah/internal/server/controller"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func InitClientRouter(api *gin.RouterGroup, i do.Injector) {
	clientGroup := api.Group("client")
	clientController := do.MustInvoke[controller.ClientController](i)
	fileController := do.MustInvoke[controller.FileController](i)

	clientGroup.GET("/:id", clientController.GetClient)
	clientGroup.GET("/page", clientController.GetClientPage)
	clientGroup.DELETE("/:id", clientController.DeleteClient)
	clientGroup.GET("/:id/pty", clientController.OpenPty)
	clientGroup.GET("/:id/systemInfo", clientController.GetClientStat)
	clientGroup.GET("/install_script", clientController.GetInstallScript)

	clientGroup.GET("/:id/process", clientController.GetClientProcessList)
	clientGroup.DELETE("/:id/process/:pid", clientController.KillClientProcess)
	clientGroup.GET("/:id/network", clientController.GetClientNetworkList)
	clientGroup.GET("/:id/docker/container", clientController.GetClientDockerContainerList)

	fileGroup := clientGroup.Group("/:id/file")
	fileGroup.GET("", fileController.GetFileList)
	fileGroup.GET("/content", fileController.GetFileContent)
	fileGroup.POST("/rename", fileController.RenameFile)
	fileGroup.DELETE("", fileController.DeleteFile)
	fileGroup.PUT("/content", fileController.UpdateFileContent)
	fileGroup.POST("", fileController.UploadFile)
	fileGroup.POST("/dir", fileController.NewDir)

}
