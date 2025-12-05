package controller

import (
	"noah/internal/service"
	"noah/pkg/errcode"
	"noah/pkg/request"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

func newAuthController() *AuthController {
	return &AuthController{
		authService: service.GetAuthService(),
	}
}

// Login 用户登录
// @Summary 用户登录
// @Description 支持手机验证码、密码登录
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginRequest body request.LoginRequest true "登录请求"
// @Success 200 {object} response.LoginResponse "成功响应"
// @Router /auth/login [post]
func (h *AuthController) Login(ctx *gin.Context) {
	var req request.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		FailWithMsg(ctx, errcode.ErrInvalidParams, err.Error())
		return
	}
	res, err := h.authService.Login(ctx, req)
	if err != nil {
		Fail(ctx, err)
	}
	Success(ctx, res)
}
