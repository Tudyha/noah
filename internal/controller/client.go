package controller

import (
	"encoding/json"
	"fmt"
	"noah/internal/service"
	"noah/pkg/config"
	"noah/pkg/errcode"
	"noah/pkg/request"
	"noah/pkg/response"
	"noah/pkg/utils"

	"github.com/gin-gonic/gin"
)

type ClientController struct {
	clientService service.ClientService
	workService   service.WorkService
}

func newClientController() *ClientController {
	return &ClientController{
		clientService: service.GetClientService(),
		workService:   service.GetWorkService(),
	}
}

func (h *ClientController) GetClientPage(ctx *gin.Context) {
	appID := GetAppID(ctx)
	if appID == 0 {
		Fail(ctx, errcode.ErrInvalidParams)
		return
	}

	var req request.ClientQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		FailWithMsg(ctx, errcode.ErrInvalidParams, err.Error())
		return
	}
	res, err := h.clientService.GetPage(ctx, appID, req)
	if err != nil {
		Fail(ctx, err)
		return
	}
	Success(ctx, res)
}

func (h *ClientController) GetClientBind(ctx *gin.Context) {
	appID := GetAppID(ctx)
	if appID == 0 {
		Fail(ctx, errcode.ErrInvalidParams)
		return
	}

	app, err := h.workService.GetAppByAppID(ctx, appID)
	if err != nil {
		Fail(ctx, err)
		return
	}

	cfg := config.Get()

	clientConfig := config.ClientConfig{
		ServerAddr: cfg.Server.TCP.Addr,
		AppId:      appID,
		AppSecret:  app.Secret,
	}

	data, err := json.Marshal(clientConfig)
	if err != nil {
		Fail(ctx, err)
		return
	}

	c := utils.Base64Encode(data)

	sh := "curl -s http://127.0.0.1:8080/file/%s -o /tmp/noah-cli && chmod +x /tmp/noah-cli && /tmp/noah-cli run -c %s"
	Success(ctx, response.ClientBindResponse{
		MacBind: fmt.Sprintf(sh, "noah-mac", c),
	})

}
