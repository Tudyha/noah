package dao

import (
	"noah/internal/server/enum"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type ChannelDao struct {
	Db *gorm.DB
}

func NewChannelDao(i do.Injector) (*ChannelDao, error) {
	return &ChannelDao{
		Db: do.MustInvoke[*gorm.DB](i),
	}, nil
}

type Channel struct {
	gorm.Model
	ChannelType enum.ChannelType   // 通道类型
	ClientId    uint               // 客户端id
	ServerPort  int                // 服务端端口
	ClientIp    string             // 客户端ip
	ClientPort  int                // 客户端端口
	Status      enum.ChannelStatus // 服务端状态
	FailReason  string
}

func (Channel) TableName() string {
	return "channel"
}

func (d ChannelDao) Save(channel Channel) (id uint, err error) {
	result := d.Db.Create(&channel)
	if result.Error != nil {
		return 0, result.Error
	}
	return channel.ID, nil
}

func (d ChannelDao) GetById(channelId uint) (channel Channel, err error) {
	err = d.Db.Where("id = ?", channelId).First(&channel).Error
	return channel, err
}

func (d ChannelDao) Delete(id uint) error {
	d.Db.Unscoped().Delete(&Channel{}, id)
	return nil
}

func (d ChannelDao) List(clientId uint) (channels []Channel, err error) {
	if clientId == 0 {
		err = d.Db.Find(&channels).Error
		return channels, err
	}
	err = d.Db.Where("client_id = ?", clientId).Find(&channels).Error
	return channels, err
}

func (d ChannelDao) UpdateStatus(id uint, status enum.ChannelStatus, failReason string) error {
	return d.Db.Model(&Channel{}).Where("id = ?", id).Updates(Channel{Status: status, FailReason: failReason}).Error
}
