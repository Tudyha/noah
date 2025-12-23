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

		token := getToken(ctx)
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusOK, response.Response{
				Code: errcode.UNAUTHORIZED,
				Msg:  "invalid token",
			})
			return
		}

		userID, err := authService.ValidateToken(ctx, token)
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

func getToken(ctx *gin.Context) string {
	token := ctx.Query(constant.HttpHeaderTokenKey)
	if token != "" {
		return token
	}
	header := ctx.GetHeader(constant.HttpHeaderTokenKey)
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}
