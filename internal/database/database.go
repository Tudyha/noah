package database

import (
	"fmt"
	"sync"
	"time"

	"noah/internal/model"
	"noah/pkg/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db   *gorm.DB
	once sync.Once
)

// Init 初始化数据库连接
func Init(cfg *config.DatabaseConfig) error {
	var err error
	once.Do(func() {
		var dialector gorm.Dialector

		switch cfg.Type {
		case "sqlite":
			dialector = sqlite.Open(cfg.DNS)
		default:
			err = fmt.Errorf("unsupported database type: %s", cfg.Type)
			return
		}

		// 配置GORM
		gormConfig := &gorm.Config{}
		switch cfg.LogLevel {
		case "info":
			gormConfig.Logger = logger.Default.LogMode(logger.Info)
		case "error":
			gormConfig.Logger = logger.Default.LogMode(logger.Error)
		default:
			gormConfig.Logger = logger.Default.LogMode(logger.Silent)
		}

		db, err = gorm.Open(dialector, gormConfig)
		if err != nil {
			return
		}

		// 获取底层数据库连接
		sqlDB, err := db.DB()
		if err != nil {
			return
		}

		// 设置连接池参数
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

		db.AutoMigrate(&model.User{}, &model.WorkSpace{}, &model.WorkSpaceUser{}, &model.WorkSpaceApp{}, &model.Client{})
	})

	return err
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return db
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
