package controller

import (
	"errors"
	"fmt"
	"github.com/golang-module/carbon/v2"
	"github.com/jinzhu/copier"
	"net/http"
	"noah/internal/server/config"
	"noah/internal/server/dto"
	"noah/internal/server/service"
	"noah/internal/server/vo"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ClientController struct{}

func NewClientController() *ClientController {
	return &ClientController{}
}

// CreateClient 新增客户端
func (c ClientController) CreateClient(ctx *gin.Context) {
	var body vo.ClientPostReq
	if err := ctx.ShouldBindJSON(&body); err != nil {
		Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var ClientPostDto dto.ClientPostDto
	copier.Copy(&ClientPostDto, &body)

	id, err := service.GetClientService().Save(ClientPostDto)
	if err != nil {
		Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	Success(ctx, id)
}

// GetClient 获取客户端列表
func (c ClientController) GetClient(ctx *gin.Context) {
	var req vo.ClientListQueryReq
	if err := ctx.BindQuery(&req); err != nil {
		Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	total, Clients := service.GetClientService().GetClient(req)

	Success(ctx, &dto.PageInfo{
		List:  Clients,
		Total: total,
	})
}

// DeleteClient 删除客户端
func (c ClientController) DeleteClient(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	//发送命令让客户端退出
	_, err := service.GetClientService().SendCommand(uint(id), "exit", "")
	if err != nil {
		Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	//断开ws连接
	err = service.GetClientService().Exit(uint(id))
	if err != nil {
		Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	//删除客户端
	err = service.GetClientService().Delete(uint(id))
	if err != nil {
		Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	Success(ctx, nil)
}

// NewWsClient 新建客户端ws连接
func (c ClientController) NewWsClient(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	ws, err := config.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println("upgrade:", err)
		Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	err = service.GetClientService().AddConnection(uint(id), ws)
	if err != nil {
		Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	//Success(ctx, nil)
}

// SendCommandHandler 发送命令
func (c ClientController) SendCommandHandler(ctx *gin.Context) {
	var form vo.SendCommandReq
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if len(strings.TrimSpace(form.Command)) == 0 {
		Fail(ctx, http.StatusBadRequest, "command is empty")
		return
	}

	id := form.ID

	res, err := service.GetClientService().SendCommand(id, form.Command, form.Parameter)
	if err != nil {
		Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	Success(ctx, res)
}

// Generate 生成客户端文件
func (c ClientController) Generate(ctx *gin.Context) {
	var req vo.ClientGenerateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	filename, err := generate(req)
	if err != nil {
		Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	Success(ctx, filename)
}

// generate 生成客户端文件
// return filename 客户端文件名
func generate(req vo.ClientGenerateReq) (string, error) {
	if len(strings.TrimSpace(req.ServerAddr)) == 0 {
		return "", errors.New("serverAddr is empty")
	}

	if len(strings.TrimSpace(req.Port)) == 0 {
		return "", errors.New("port is empty")
	}

	filename, err := service.GetClientService().Generate(req.ServerAddr, req.Port, req.OsType, req.Filename)
	if err != nil {
		return "", err
	}
	return filename, nil
}

// Update 更新客户端
func (c ClientController) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	//生成最新客户端
	var req vo.ClientGenerateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	filename, err := generate(req)
	if err != nil {
		Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	//发送命令让客户端升级
	_, err = service.GetClientService().SendCommand(uint(id), "update", filename)
	if err != nil {
		Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	Success(ctx, "success")
}

func (c ClientController) GetClientInfo(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	start := ctx.Query("start")
	end := ctx.Query("end")
	var startTime, endTime time.Time
	if start == "" || end == "" {
		//获取当前时间
		endTime = time.Now()
		//获取5分钟前时间
		startTime = endTime.Add(-5 * time.Minute)
	} else {
		//15:39:55 转time.Time
		startTime = carbon.Parse(start).StdTime()
		endTime = carbon.Parse(end).StdTime()
	}

	clientInfoList, err := service.GetClientService().GetSystemInfo(uint(id), startTime, endTime)
	if err != nil {
		Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	Success(ctx, clientInfoList)
}
