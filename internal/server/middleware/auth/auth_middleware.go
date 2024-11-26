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
	tempTokenMap    sync.Map
}

func NewAuthMiddleware(i do.Injector) (*AuthMiddleware, error) {
	a := AuthMiddleware{
		tokenMap:        sync.Map{},
		refreshTokenMap: sync.Map{},
		userService:     do.MustInvoke[service.IUserService](i),
		tempTokenMap:    sync.Map{},
	}

	a.loadFromDb()

	return &a, nil
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

	auth.userService.UpdateToken(id, r.Token, r.RefreshToken, r.ExpireTime, r.RefreshExpireTime)

	return r
}

func (auth *AuthMiddleware) loadFromDb() {
	users, err := auth.userService.GetUserList()
	if err != nil {
		return
	}
	for _, user := range users {
		if user.Token != "" && user.ExpireTime.After(time.Now()) {
			r := response.AuthResult{
				UserId:            user.ID,
				Token:             user.Token,
				RefreshToken:      user.RefreshToken,
				ExpireTime:        user.ExpireTime,
				RefreshExpireTime: user.RefreshExpireTime,
			}
			auth.tokenMap.Store(user.Token, r)
			auth.refreshTokenMap.Store(user.RefreshToken, r)
		} else {
			if user.RefreshToken != "" && user.RefreshExpireTime.After(time.Now()) {
				auth.GenerateToken(user.ID)
			}
		}
	}
}

func (auth *AuthMiddleware) GenerateTempToken() response.AuthResult {
	token := utils.RandString(16)
	r := response.AuthResult{
		UserId:     0,
		Token:      token,
		ExpireTime: time.Now().Add(expireTime),
	}
	auth.tempTokenMap.Store(token, r)
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

func (auth *AuthMiddleware) TempAuthMiddlewareFunc() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := getToken(c)
		if token == "" {
			authFail(c)
			return
		}
		if data, ok := auth.tempTokenMap.Load(token); !ok {
			authFail(c)
		} else {
			if data.(response.AuthResult).ExpireTime.Before(time.Now()) {
				authFail(c)
				auth.tempTokenMap.Delete(token)
			} else {
				c.Next()
			}
		}
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
