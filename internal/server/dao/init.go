package dao

import (
	"noah/internal/server/model"
	"os"

	"github.com/samber/do/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init(i do.Injector) {
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
	do.Provide(i, NewUserDao)
	do.Provide(i, NewClientDao)
	do.Provide(i, NewClientStatDao)
	do.Provide(i, NewTunnelDao)

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Client{})
	db.AutoMigrate(&model.ClientStat{})
	db.AutoMigrate(&model.Tunnel{})
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

func Paginate(page int, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page < 1 {
			page = 1
		}

		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}
