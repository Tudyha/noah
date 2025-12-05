package logger

import (
	"noah/pkg/config"
	"sync"
)

// Logger 定义了应用层使用的基本日志接口
type Logger interface {
	// Debug 记录调试级别的日志
	Debug(msg string, fields ...interface{})

	// Info 记录信息级别的日志（通常用于正常操作）
	Info(msg string, fields ...interface{})

	// Warn 记录警告级别的日志
	Warn(msg string, fields ...interface{})

	// Error 记录错误级别的日志
	Error(msg string, fields ...interface{})

	// Sync 同步日志缓冲区，确保所有日志都被写出
	Sync() error
}

type Field struct {
	Key   string
	Value interface{}
}

var (
	log  Logger
	once sync.Once
)

// Init 初始化日志
func Init(cfg *config.LoggerConfig) error {
	var err error
	once.Do(func() {
		log, err = initZap(cfg)
	})
	return err
}

// GetLogger 获取日志器
func GetLogger() Logger {
	return log
}

// Sync 同步日志缓冲区
func Sync() error {
	if log != nil {
		return log.Sync()
	}
	return nil
}

func Debug(msg string, fields ...interface{}) {
	log.Debug(msg, fields...)
}

func Info(msg string, fields ...interface{}) {
	log.Info(msg, fields...)
}

func Warn(msg string, fields ...interface{}) {
	log.Warn(msg, fields...)
}

func Error(msg string, fields ...interface{}) {
	log.Error(msg, fields...)
}
