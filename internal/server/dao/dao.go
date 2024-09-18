package dao

import (
	"noah/internal/server/environment"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db            *gorm.DB
	clientDao     *ClientDao
	clientInfoDao *ClientInfoDao
	channelDao    *ChannelDao
)

func InitDb(dbConfig environment.DatabaseConfig) (err error) {
	//db, err = openMysqlDb(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	//	dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName))
	db, err = openSqlLiteDb()
	if err != nil {
		panic(err)
	}
	if db.Error != nil {
		panic(db.Error)
	}

	_ = db.AutoMigrate(&Client{})
	_ = db.AutoMigrate(&ClientInfo{})
	_ = db.AutoMigrate(&Channel{})

	clientDao = &ClientDao{db}
	clientInfoDao = &ClientInfoDao{db}
	channelDao = &ChannelDao{db}

	return nil
}

func openMysqlDb(connectStr string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(connectStr), &gorm.Config{})
}

func openSqlLiteDb() (*gorm.DB, error) {
	// 确保 data 目录存在
	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		return nil, err
	}
	return gorm.Open(sqlite.Open("data/noah.db"), &gorm.Config{})
}

func GetClientDao() *ClientDao {
	return clientDao
}

func GetClientInfoDao() *ClientInfoDao {
	return clientInfoDao
}

func GetChannelDao() *ChannelDao {
	return channelDao
}
