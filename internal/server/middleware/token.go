package middleware

import (
	"errors"
	"fmt"
	"noah/internal/server/dao"
	"noah/internal/server/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muesli/cache2go"
	"github.com/samber/do/v2"
)

var (
	tokenCache        = cache2go.Cache("token")
	refreshTokenCache = cache2go.Cache("refreshToken")
	expireTime        = time.Second * 60 * 3
	refreshTime       = time.Second * 60 * 60 * 24 * 7
)

type AuthMiddleware struct {
	tokenDao *dao.TokenDao
}

func NewAuthMiddleware(i do.Injector) *AuthMiddleware {
	authMiddleware := &AuthMiddleware{
		tokenDao: do.MustInvoke[*dao.TokenDao](i),
	}
	authMiddleware.initFromDb()

	tokenCache.SetAddedItemCallback(func(entry *cache2go.CacheItem) {
		fmt.Println("tokenCache", entry.Data().(AuthResult))
		authMiddleware.saveToDb(entry.Data().(AuthResult))
	})
	tokenCache.SetAboutToDeleteItemCallback(func(entry *cache2go.CacheItem) {
		authMiddleware.delDb(entry.Data().(AuthResult).Token)
	})
	return authMiddleware
}

type AuthResult struct {
	AuthDomain        authDomain
	Token             string    `json:"token"`
	RefreshToken      string    `json:"refreshToken"`
	ExpireTime        time.Time `json:"expireTime"`
	RefreshExpireTime time.Time `json:"refreshExpireTime"`
}

type authDomain uint8

const (
	authDomainAdmin authDomain = iota + 1
	authDomainClient
)

func (auth *AuthMiddleware) initFromDb() {
	list, err := auth.tokenDao.List()
	if err != nil {
		return
	}
	for _, token := range list {
		a := AuthResult{
			AuthDomain:        authDomain(token.AuthDomain),
			Token:             token.Token,
			RefreshToken:      token.RefreshToken,
			ExpireTime:        token.ExpireTime,
			RefreshExpireTime: token.RefreshExpireTime,
		}
		if !token.ExpireTime.Before(time.Now()) {
			tokenCache.Add(token.Token, time.Until(token.ExpireTime), a)
		}
		if !token.RefreshExpireTime.Before(time.Now()) {
			refreshTokenCache.Add(token.RefreshToken, time.Until(token.RefreshExpireTime), a)
		}
	}
}

func (auth *AuthMiddleware) saveToDb(a AuthResult) {
	auth.tokenDao.Save(dao.Token{
		AuthDomain:        uint8(a.AuthDomain),
		Token:             a.Token,
		RefreshToken:      a.RefreshToken,
		ExpireTime:        a.ExpireTime,
		RefreshExpireTime: a.RefreshExpireTime,
	})
}

func (auth *AuthMiddleware) delDb(token string) {
	auth.tokenDao.Delete(token)
}

func GenerateAdminToken() AuthResult {
	return generateToken(authDomainAdmin)
}

func GenerateClientToken() AuthResult {
	return generateToken(authDomainClient)
}

func generateToken(authDomain authDomain) AuthResult {
	token := utils.RandString(16)
	refreshToken := utils.RandString(16)
	r := AuthResult{
		AuthDomain:        authDomain,
		Token:             token,
		RefreshToken:      refreshToken,
		ExpireTime:        time.Now().Add(expireTime),
		RefreshExpireTime: time.Now().Add(refreshTime),
	}
	tokenCache.Add(token, expireTime, r)
	refreshTokenCache.Add(refreshToken, refreshTime, r)
	return r
}

func RefreshClientToken(refreshToken string) (AuthResult, error) {
	if auth, err := refreshTokenCache.Value(refreshToken); err != nil {
		return AuthResult{}, errors.New("refresh token is invalid")
	} else {
		if auth.Data().(AuthResult).AuthDomain != authDomainClient {
			return AuthResult{}, errors.New("refresh token is invalid")
		}
		return generateToken(auth.Data().(AuthResult).AuthDomain), nil
	}
}

func getToken(c *gin.Context) string {
	token := c.GetHeader("Token")
	if token == "" {
		token = c.Query("token")
	}
	return token
}

func AdminAuthMiddlewareFunc() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := getToken(c)
		if token == "" {
			authFail(c)
			return
		}
		if auth, err := tokenCache.Value(token); err != nil {
			authFail(c)
		} else {
			if auth.Data().(AuthResult).AuthDomain != authDomainAdmin {
				authFail(c)
			} else if auth.Data().(AuthResult).ExpireTime.Before(time.Now()) {
				authFail(c)
				tokenCache.Delete(token)
			}
		}
	}
}

func authFail(c *gin.Context) {
	c.JSON(401, gin.H{
		"code": 401,
		"msg":  "token is invalid",
	})
	c.Abort()
}

func ClientAuthMiddlewareFunc() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := getToken(c)
		if token == "" {
			authFail(c)
			return
		}
		tokenCache.Foreach(func(key interface{}, value *cache2go.CacheItem) {
			fmt.Println(key, value.Data().(AuthResult), value.LifeSpan())
		})
		if auth, err := tokenCache.Value(token); err != nil {
			authFail(c)
		} else {
			if auth.Data().(AuthResult).AuthDomain != authDomainClient {
				authFail(c)
			} else if auth.Data().(AuthResult).ExpireTime.Before(time.Now()) {
				authFail(c)
				tokenCache.Delete(token)
			}
		}
	}
}

func AuthMiddlewareFunc() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := getToken(c)
		if token == "" {
			authFail(c)
			return
		}
		if auth, err := tokenCache.Value(token); err != nil {
			authFail(c)
		} else if auth.Data().(AuthResult).ExpireTime.Before(time.Now()) {
			authFail(c)
			tokenCache.Delete(token)
		}
	}
}
