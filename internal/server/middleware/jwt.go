package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	"noah/internal/server/dto"
	"time"
)

var (
	realm         = "noah"
	identityKey   = "userId"
	secretKey     = "noah"
	adminPassword = ""
	auth          *jwt.GinJWTMiddleware
)

func SetAdminPassword(password string) {
	adminPassword = password
}

func RegisterJwtMiddleWare() (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(initParams())
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
		return nil, err
	}

	auth = authMiddleware
	return authMiddleware, nil
}

func initParams() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:       realm,
		Key:         []byte(secretKey),
		Timeout:     time.Hour * 24,
		MaxRefresh:  time.Hour * 24 * 7,
		IdentityKey: identityKey,
		PayloadFunc: payloadFunc(),

		IdentityHandler: identityHandler(),
		Authenticator:   authenticator(),
		Unauthorized:    unauthorized(),
		LoginResponse:   loginResponse(),
		TokenLookup:     "header: Authorization, query:token",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
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

		if username == "admin" && password == adminPassword {
			return &dto.UserInfo{
				UserID:   1,
				Username: username,
			}, nil
		}
		return nil, jwt.ErrFailedAuthentication
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

func GetToken() (string, error) {
	token, _, err := auth.TokenGenerator(&dto.UserInfo{
		UserID:   1,
		Username: "admin",
	})
	if err != nil {
		return "", err
	}
	return token, nil
}
