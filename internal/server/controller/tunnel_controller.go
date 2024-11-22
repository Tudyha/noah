package controller

import (
	"noah/internal/server/service"
	"noah/pkg/errcode"
	"noah/pkg/request"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type TunnelController struct {
	tunnelService service.ITunnelService
}

func NewTunnelController(i do.Injector) (TunnelController, error) {
	return TunnelController{
		tunnelService: do.MustInvoke[service.ITunnelService](i),
	}, nil
}

func (c TunnelController) NewTunnel(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	uintId := uint(id)

	var req request.CreateTunnelReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		Fail(ctx, errcode.ErrInvalidParameter)
		return
	}

	err = c.tunnelService.NewTunnel(uintId, req)
	if err != nil {
		Fail(ctx, errcode.ErrInternalError)
		return
	}

	Success(ctx, "success")
}

func (c TunnelController) GetTunnelList(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	uintId := uint(id)

	res, err := c.tunnelService.GetTunnelList(uintId)
	if err != nil {
		Fail(ctx, errcode.ErrInternalError)
		return
	}
	Success(ctx, res)

}

func (c TunnelController) DeleteTunnel(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("tunnelId"))
	uintId := uint(id)

	err := c.tunnelService.DeleteTunnel(uintId)
	if err != nil {
		Fail(ctx, errcode.ErrInternalError)
		return
	}
	Success(ctx, "success")
}
