package middleware

import (
	"net/http"
	"noah/internal/service"
	"noah/pkg/constant"
	"noah/pkg/errcode"
	"noah/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

func Auth(authService service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uri := ctx.Request.RequestURI
		if lo.Contains(constant.White_Api_List, uri) {
			ctx.Next()
			return
		}

		header := ctx.GetHeader(constant.HttpHeaderTokenKey)
		if header == "" {
			ctx.AbortWithStatusJSON(http.StatusOK, response.Response{
				Code: errcode.UNAUTHORIZED,
				Msg:  "token is empty",
			})
			return
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusOK, response.Response{
				Code: errcode.UNAUTHORIZED,
				Msg:  "invalid token format",
			})
			return
		}
		userID, err := authService.ValidateToken(ctx, parts[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusOK, response.Response{
				Code: errcode.UNAUTHORIZED,
				Msg:  "invalid token",
			})
			return
		}
		ctx.Set(constant.HttpHeaderUserIDKey, userID)
		ctx.Next()
	}
}
