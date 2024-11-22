package controller

import (
	"noah/internal/server/middleware/auth"
	"noah/internal/server/service"
	"noah/pkg/errcode"
	"noah/pkg/request"
	"noah/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type UserController struct {
	userService    service.IUserService
	authMiddleware *auth.AuthMiddleware
}

func NewUserController(i do.Injector) (UserController, error) {
	return UserController{
		userService:    do.MustInvoke[service.IUserService](i),
		authMiddleware: do.MustInvoke[*auth.AuthMiddleware](i),
	}, nil
}

func (c UserController) Login(ctx *gin.Context) {
	var req request.LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Fail(ctx, errcode.ErrInvalidParameter)
		return
	}

	user, err := c.userService.Login(req.Username, req.Password)
	if err != nil {
		Fail(ctx, err)
		return
	}

	authResult := c.authMiddleware.GenerateToken(user.ID)

	err = c.userService.UpdateToken(user.ID, authResult.Token, authResult.RefreshToken, authResult.ExpireTime, authResult.RefreshExpireTime)
	if err != nil {
		Fail(ctx, err)
		return
	}

	Success(ctx, response.LoginResp{
		AuthResult: authResult,
	})
}

func (c UserController) GetUser(ctx *gin.Context) {
	userId, ok := ctx.Get("userId")
	if !ok {
		Fail(ctx, errcode.ErrTokenExpired)
		return
	}
	user, err := c.userService.GetUser(userId.(uint))
	if err != nil {
		Fail(ctx, err)
		return
	}
	Success(ctx, response.GetUserResp{
		UserResp: response.UserResp{
			Avatar:       user.Avatar,
			Introduction: user.Introduction,
			Name:         user.Name,
			Roles:        []string{"admin"},
			UserId:       user.ID,
		},
	})
}

func (c UserController) GetUserPage(ctx *gin.Context) {
	var req request.PageQuery
	if err := ctx.ShouldBindQuery(&req); err != nil {
		Fail(ctx, errcode.ErrInvalidParameter)
		return
	}

	total, users, err := c.userService.GetUserPage(req.Page, req.Size)
	if err != nil {
		Fail(ctx, err)
		return
	}

	res := response.GetUserPageResp{
		Total: total,
	}

	for _, user := range users {
		res.List = append(res.List, response.UserResp{
			Avatar:       user.Avatar,
			Introduction: user.Introduction,
			LoginTime:    user.LoginTime,
			Name:         user.Name,
			UserId:       user.ID,
		})
	}

	Success(ctx, res)
}

func (c UserController) UpdatePassword(ctx *gin.Context) {
	currentUserId, ok := ctx.Get("userId")
	if !ok {
		Fail(ctx, errcode.ErrTokenExpired)
		return
	}
	if currentUserId.(uint) != uint(1) {
		Fail(ctx, errcode.ErrPermissionDenied)
		return
	}

	var req request.UpdatePasswordReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Fail(ctx, errcode.ErrInvalidParameter)
		return
	}

	userId, err := strconv.Atoi(ctx.Param("userId"))
	if err != nil {
		Fail(ctx, errcode.ErrTokenExpired)
		return
	}
	err = c.userService.UpdatePassword(uint(userId), req.Password)
	if err != nil {
		Fail(ctx, err)
		return
	}
	Success(ctx, nil)
}
