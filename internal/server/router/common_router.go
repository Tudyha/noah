package router

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func InitCommonRouter(api *gin.RouterGroup, i do.Injector) {
	api.GET("/file/download/:filename", func(c *gin.Context) {
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
}
