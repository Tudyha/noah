package dao

import (
	"noah/internal/server/environment"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db       *gorm.DB
	DeviceDa *DeviceDao
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

	err = db.AutoMigrate(&Device{})
	if err != nil {
		//log error
	}

	DeviceDa = &DeviceDao{db}

	return nil
}

func openMysqlDb(connectStr string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(connectStr), &gorm.Config{})
}

func openSqlLiteDb() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("data/noah.db"), &gorm.Config{})
}
