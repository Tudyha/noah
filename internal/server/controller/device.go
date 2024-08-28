package controller

import (
	"net/http"
	"noah/internal/server/dto"
	"noah/internal/server/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type DeviceController struct{}

func NewDeviceController() *DeviceController {
	return &DeviceController{}
}

type DevicePostVo struct {
	Hostname    string `json:"hostname"`
	Username    string `json:"username"`
	UserID      string `json:"userId"`
	OSName      string `json:"osName"`
	OSArch      string `json:"osArch"`
	MacAddress  string `json:"macAddress"`
	IPAddress   string `json:"ipAddress"`
	Port        string `json:"port"`
	FetchedUnix int64  `json:"fetchedUnix"`
}

func (d *DeviceController) CreateDevice(c *gin.Context) {
	var body DevicePostVo
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	var devicePostDto dto.DevicePostDto
	copier.Copy(&devicePostDto, &body)

	id, err := service.GetDeviceService().Save(devicePostDto)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	Success(c, id)
}

func (d *DeviceController) GetDevice(c *gin.Context) {
	var req dto.DeviceListQueryDto
	if err := c.BindQuery(&req); err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	total, devices := service.GetDeviceService().GetDevice(req)

	Success(c, &dto.PageInfo{
		List: devices,
		// HasNextPage: false,
		// NextPage:    2,
		// PageNum:     1,
		// PageSize:    100,
		Total: total,
		// Pages:       1,
		// Size:        1,
	})
}

func (d *DeviceController) DeleteDevice(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := service.GetClientService().SendCommand(uint(id), "exit", "")
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	err = service.GetClientService().Exit(uint(id))
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	err = service.GetDeviceService().Delete(uint(id))
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, nil)
}
