package dao

import (
	"os"

	"github.com/samber/do/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init(i do.Injector) (err error) {
	//db, err = openMysqlDb(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	//	dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName))
	db, err := openSqlLiteDb()
	if err != nil {
		panic(err)
	}
	if db.Error != nil {
		panic(db.Error)
	}

	do.ProvideValue(i, db)

	_ = db.AutoMigrate(&Client{})
	_ = db.AutoMigrate(&ClientInfo{})
	_ = db.AutoMigrate(&Channel{})
	_ = db.AutoMigrate(&Token{})

	do.Provide(i, NewClientDao)
	do.Provide(i, NewClientInfoDao)
	do.Provide(i, NewChannelDao)
	do.Provide(i, NewTokenDao)

	return nil
}

// func openMysqlDb(connectStr string) (*gorm.DB, error) {
// 	return gorm.Open(mysql.Open(connectStr), &gorm.Config{})
// }

func openSqlLiteDb() (*gorm.DB, error) {
	// 确保 data 目录存在
	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		return nil, err
	}
	return gorm.Open(sqlite.Open("data/noah.db"), &gorm.Config{})
}
