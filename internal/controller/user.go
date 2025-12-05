package controller

import (
	"noah/internal/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.UserService
}

func newUserController() *UserController {
	return &UserController{
		userService: service.GetUserService(),
	}
}

// GetUser 获取用户信息
// @Summary 获取用户信息
// @Description 获取用户信息
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} response.UserResponse "成功响应"
// @Router /user [get]
func (h *UserController) GetUser(ctx *gin.Context) {
	userId := GetUserId(ctx)
	user, err := h.userService.GetByID(ctx, userId)
	if err != nil {
		Fail(ctx, err)
		return
	}
	Success(ctx, user)
}
