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
		return
	}
	Success(ctx, res)
}

// Register 用户注册
// @Summary 用户注册
// @Description 支持邮箱注册
// @Tags Auth
// @Accept json
// @Produce json
// @Param registerRequest body request.RegisterRequest true "注册请求"
// @Success 200 {object} response.Response "成功响应"
// @Router /auth/register [post]
func (h *AuthController) Register(ctx *gin.Context) {
	var req request.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		FailWithMsg(ctx, errcode.ErrInvalidParams, err.Error())
		return
	}
	err := h.authService.Register(ctx, req)
	if err != nil {
		Fail(ctx, err)
		return
	}
	Success(ctx, nil)
}

// SendCode 发送验证码
// @Summary 发送验证码
// @Description 支持手机号、邮箱发送验证码
// @Tags Auth
// @Accept json
// @Produce json
// @Param sendCodeRequest body request.SendCodeRequest true "发送验证码请求"
// @Success 200 {object} response.Response "成功响应"
// @Router /auth/send_code [post]
func (h *AuthController) SendCode(ctx *gin.Context) {
	var req request.SendCodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		FailWithMsg(ctx, errcode.ErrInvalidParams, err.Error())
		return
	}
	err := h.authService.SendCode(ctx, req)
	if err != nil {
		Fail(ctx, err)
		return
	}
	Success(ctx, nil)
}
