package logger

import (
	"noah/pkg/config"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zapLogger struct {
	sugar *zap.SugaredLogger
	log   *zap.Logger
}

func (z *zapLogger) Debug(msg string, fields ...interface{}) {
	z.sugar.Debugw(msg, fields...)
}

func (z *zapLogger) Info(msg string, fields ...interface{}) {
	z.sugar.Infow(msg, fields...)
}

func (z *zapLogger) Warn(msg string, fields ...interface{}) {
	z.sugar.Warnw(msg, fields...)
}

func (z *zapLogger) Error(msg string, fields ...interface{}) {
	z.sugar.Errorw(msg, fields...)
}

func (z *zapLogger) Sync() error {
	return z.log.Sync()
}

func initZap(cfg *config.LoggerConfig) (l Logger, err error) {
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 输出配置
	var writeSyncer zapcore.WriteSyncer
	if cfg.Output == "file" {
		// 确保日志目录存在
		logDir := filepath.Dir(cfg.FilePath)
		if err = os.MkdirAll(logDir, 0755); err != nil {
			return
		}

		// 文件输出
		writeSyncer = zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize,    // MB
			MaxBackups: cfg.MaxBackups, // 保留旧文件的最大个数
			MaxAge:     cfg.MaxAge,     // days
			Compress:   cfg.Compress,   // 是否压缩
		})
	} else {
		// 控制台输出
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// 创建编码器
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 创建核心
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// 创建日志器
	zlog := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
	sugar := zlog.Sugar()

	return &zapLogger{
		sugar: sugar,
		log:   zlog,
	}, nil
}
