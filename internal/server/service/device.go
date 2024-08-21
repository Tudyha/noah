package service

import "noah/internal/server/dto"

var (
	deviceServiceInstance IDeviceService
)

func RegisterDeviceService(i IDeviceService) {
	deviceServiceInstance = i
}

func GetDeviceService() IDeviceService {
	return deviceServiceInstance
}

type IDeviceService interface {
	Save(body dto.DevicePostDto) (id uint, err error)
	UpdateStatus(id uint, status int8)
	GetDevice(query dto.DeviceListQueryDto) (total int64, deviceDtos []dto.DeviceDto)
	ScheduleUpdateStatus() error
}
