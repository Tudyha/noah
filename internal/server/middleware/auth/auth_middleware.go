package auth

import (
	"errors"
	"noah/internal/server/service"
	"noah/pkg/errcode"
	"noah/pkg/response"
	"noah/pkg/utils"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

var (
	expireTime  = time.Second * 60 * 60 * 24 * 1
	refreshTime = time.Second * 60 * 60 * 24 * 7
)

type AuthMiddleware struct {
	tokenMap        sync.Map
	refreshTokenMap sync.Map
	userService     service.IUserService
}

func NewAuthMiddleware(i do.Injector) (AuthMiddleware, error) {
	return AuthMiddleware{
		tokenMap:        sync.Map{},
		refreshTokenMap: sync.Map{},
		userService:     do.MustInvoke[service.IUserService](i),
	}, nil
}

func (auth *AuthMiddleware) GenerateToken(id uint) response.AuthResult {
	token := utils.RandString(16)
	refreshToken := utils.RandString(16)
	r := response.AuthResult{
		UserId:            id,
		Token:             token,
		RefreshToken:      refreshToken,
		ExpireTime:        time.Now().Add(expireTime),
		RefreshExpireTime: time.Now().Add(refreshTime),
	}
	auth.tokenMap.Store(token, r)
	auth.refreshTokenMap.Store(refreshToken, r)
	return r
}

func (auth *AuthMiddleware) RefreshToken(refreshToken string) (response.AuthResult, error) {
	if data, ok := auth.refreshTokenMap.Load(refreshToken); !ok {
		return response.AuthResult{}, errors.New("refresh token is invalid")
	} else {
		if data.(response.AuthResult).RefreshExpireTime.Before(time.Now()) {
			return response.AuthResult{}, errors.New("refresh token is invalid")
		}

		return auth.GenerateToken(data.(response.AuthResult).UserId), nil
	}
}

func (auth *AuthMiddleware) AuthMiddlewareFunc() func(c *gin.Context) {
	return func(c *gin.Context) {
		auth.authCheck(c, getToken(c))
	}
}

func getToken(c *gin.Context) string {
	token := c.GetHeader("Token")
	if token == "" {
		token = c.Query("token")
	}
	return token
}

func (auth *AuthMiddleware) authCheck(c *gin.Context, token string) {
	if token == "" {
		authFail(c)
		return
	}
	if data, ok := auth.tokenMap.Load(token); !ok {
		authFail(c)
	} else {
		if data.(response.AuthResult).ExpireTime.Before(time.Now()) {
			authFail(c)
			auth.tokenMap.Delete(token)
		} else {
			c.Set("userId", data.(response.AuthResult).UserId)
			c.Next()
		}
	}
}

func authFail(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": errcode.ErrTokenExpired.Code,
		"msg":  errcode.ErrTokenExpired.Msg,
	})
	c.Abort()
}
