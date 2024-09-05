package log

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var logger *logrus.Logger

func SetupLogger() *logrus.Logger {
	// 创建 Logrus 日志实例
	logger = logrus.New()

	// 设置日志输出
	logger.Out = os.Stdout // 输出到标准输出，也可以是文件

	// 设置日志级别
	logger.SetLevel(logrus.InfoLevel)

	// 设置日志格式
	//logger.SetFormatter(&CustomJSONFormatter{
	//	JSONFormatter: logrus.JSONFormatter{
	//		DataKey:         "data",
	//		TimestampFormat: time.DateTime,
	//	},
	//})
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.DateTime,
		FullTimestamp:   true,
	})

	return logger
}

func Logger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始的时间
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		// 计算请求耗时
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		logger.WithFields(logrus.Fields{
			"method":    method,
			"path":      path,
			"query":     query,
			"client_ip": clientIP,
			"status":    statusCode,
			"cost_time": latency.Milliseconds(),
		}).Info("[completed request]")
	}
}

func Info(msg string, fields map[string]interface{}) {
	logger.WithFields(fields).Infoln(msg)
}

func Warn(msg string, fields map[string]interface{}) {
	logger.WithFields(fields).Warningln(msg)
}

func Error(msg string, fields map[string]interface{}) {
	logger.WithFields(fields).Errorln(msg)
}

// CustomJSONFormatter 自定义Formatter类
type CustomJSONFormatter struct {
	logrus.JSONFormatter
}

type CustomJSONFormatterFields struct {
	Time  string      `json:"time"`
	Level string      `json:"level"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

// Format 重写Format方法
func (f *CustomJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	fields := CustomJSONFormatterFields{
		Time:  entry.Time.Format(f.TimestampFormat),
		Level: entry.Level.String(),
		Msg:   entry.Message,
		Data:  entry.Data,
	}

	// 将有序的map序列化为JSON
	jsonData, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}

	// 添加换行符
	return append(jsonData, '\n'), nil
}
