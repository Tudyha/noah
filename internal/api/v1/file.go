package v1

import (
	"noah/pkg/response"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func RegisterFileRoutes(router *gin.RouterGroup) {
	router.GET("/:name", func(ctx *gin.Context) {
		filename := ctx.Params.ByName("name")
		if filename == "" {
			ctx.JSON(200, response.Response{
				Code: 400,
				Msg:  "filename is empty",
			})
			return
		}
		filename = filepath.Clean(filename)

		ctx.File("./build/" + filename)
	})
}
