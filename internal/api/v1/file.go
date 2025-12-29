package v1

import (
	"net/http"
	"noah/pkg/response"
	"noah/pkg/utils"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func RegisterFileRoutes(router *gin.RouterGroup) {
	router.GET("/:name", func(ctx *gin.Context) {
		filename := ctx.Params.ByName("name")
		if filename == "" {
			ctx.JSON(http.StatusOK, response.Response{
				Code: 400,
				Msg:  "filename is empty",
			})
			return
		}
		filename = filepath.Clean(filename)

		if utils.FileExists("./temp/" + filename) {
			ctx.File("./temp/" + filename)
			return
		}

		ctx.File("./build/" + filename)
	})
}
