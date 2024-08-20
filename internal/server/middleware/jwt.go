package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	"noah/internal/server/dto"
	"time"
)

var (
	Realm       = "noah"
	IdentityKey = "userId"
)

func RegisterJwtMiddleWare() (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(initParams())
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
		return nil, err
	}
	return authMiddleware, nil
}

func initParams() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:       Realm,
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: IdentityKey,
		PayloadFunc: payloadFunc(),

		IdentityHandler: identityHandler(),
		Authenticator:   authenticator(),
		//Authorizator:    authorizator(),
		Unauthorized:  unauthorized(),
		LoginResponse: loginResponse(),
		TokenLookup:   "header: Authorization, query:token",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}
}

func payloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*dto.UserInfo); ok {
			return jwt.MapClaims{
				"userId": v.UserID,
			}
		}
		return jwt.MapClaims{}
	}
}

func identityHandler() func(c *gin.Context) interface{} {
	return func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		return &dto.UserInfo{
			UserID: uint(claims["userId"].(float64)),
		}
	}
}

func authenticator() func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var loginVals dto.LoginReq
		if err := c.ShouldBind(&loginVals); err != nil {
			return "", jwt.ErrMissingLoginValues
		}
		username := loginVals.Username
		password := loginVals.Password

		if (username == "admin" && password == "123456") || (username == "test" && password == "test") {
			return &dto.UserInfo{
				UserID:   1,
				Username: username,
			}, nil
		}
		return nil, jwt.ErrFailedAuthentication
	}
}

func authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		if v, ok := data.(*dto.UserInfo); ok && v.Username == "admin" {
			return true
		}
		return false
	}
}

func loginResponse() func(c *gin.Context, code int, message string, time time.Time) {
	return func(c *gin.Context, code int, message string, time time.Time) {
		c.JSON(code, gin.H{
			"code":    0,
			"message": "success",
			"data":    gin.H{"token": message},
		})
	}
}

func unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		})
	}
}
