package dao

import (
	"gorm.io/gorm"
	"noah/internal/server/enum"
)

type ChannelDao struct {
	Db *gorm.DB
}

type Channel struct {
	gorm.Model
	ChannelType enum.ChannelType // 通道类型
	ClientId    uint             // 客户端id
	ServerPort  int              // 服务端端口
	ClientIp    string           // 客户端ip
	ClientPort  int              // 客户端端口
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

func (d ChannelDao) List(clientId uint) (channels []Channel, err error) {
	err = d.Db.Where("client_id = ?", clientId).Find(&channels).Error
	return channels, err
}
