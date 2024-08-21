package device

import (
	"github.com/jinzhu/copier"
	"noah/internal/server/dao"
	"noah/internal/server/dto"
	"noah/internal/server/service"
)

type deviceService struct{}

func init() {
	service.RegisterDeviceService(&deviceService{})
}

func (ds *deviceService) Save(body dto.DevicePostDto) (id uint, err error) {
	var device dao.Device
	copier.Copy(&device, &body)

	old := dao.DeviceDa.GetByMacAddress(device.MacAddress)
	if old.ID != 0 {
		// 已存在，更新数据即可
		device.ID = old.ID
		err = dao.DeviceDa.Update(device)
		if err != nil {
			return 0, err
		}
		return old.ID, nil
	}

	id, err = dao.DeviceDa.Save(device)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (ds *deviceService) UpdateStatus(id uint, status int8) {
	dao.DeviceDa.UpdateStatus(id, status)
}

func (ds *deviceService) GetDevice(query dto.DeviceListQueryDto) (total int64, deviceDtos []dto.DeviceDto) {
	total, devices := dao.DeviceDa.Page(query)

	if total == 0 {
		return 0, make([]dto.DeviceDto, 0)
	}

	copier.Copy(&deviceDtos, &devices)
	//for idx := range deviceDtos {
	//deviceDto := &deviceDtos[idx]
	// 如果上次在线时间距离当前时间超过10分钟，则认为该设备已离线
	//if deviceDto.LastOnlineTime.Add(10 * 60 * time.Second).Before(time.Now()) {
	//	deviceDto.Status = enum.DEVICE_OFFLINE
	//}
	//}
	return total, deviceDtos
}

func (ds *deviceService) ScheduleUpdateStatus() error {
	dao.DeviceDa.ScheduleUpdateStatus()
	return nil
}
