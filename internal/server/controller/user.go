package controller

import (
	"net/http"
	"noah/internal/server/environment"
	"noah/internal/server/middleware"
	"noah/internal/server/request"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type UserController struct {
	AdminPassword string
}

func NewUserController(i do.Injector) *UserController {
	return &UserController{
		AdminPassword: do.MustInvoke[environment.Environment](i).Admin.Password,
	}
}

func (u UserController) Login(c *gin.Context) {
	var loginVals request.LoginReq
	if err := c.ShouldBind(&loginVals); err != nil {
		Fail(c, http.StatusUnauthorized, err.Error())
		return
	}
	username := loginVals.Username
	password := loginVals.Password
	if username != "admin" {
		Fail(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	if password != u.AdminPassword {
		Fail(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}
	r := middleware.GenerateAdminToken()
	Success(c, r)
}

func (u UserController) RefreshToken(c *gin.Context) {
	refreshToken := c.Query("refreshToken")
	if refreshToken == "" {
		Fail(c, http.StatusUnauthorized, "refreshToken is empty")
		return
	}
	auth, err := middleware.RefreshClientToken(refreshToken)
	if err != nil {
		Fail(c, http.StatusUnauthorized, err.Error())
		return
	}
	Success(c, auth)
}
