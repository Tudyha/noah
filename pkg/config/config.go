package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var cfg *Config

// Config 配置结构体
type Config struct {
	Server      ServerConfig   `mapstructure:"server"`
	Database    DatabaseConfig `mapstructure:"database"`
	Redis       RedisConfig    `mapstructure:"redis"`
	TokenConfig TokenConfig    `mapstructure:"token"`
	Logger      LoggerConfig   `mapstructure:"logger"`
}

// ServerConfig 服务配置
type ServerConfig struct {
	Env   string      `mapstructure:"env"`
	HTTP  HTTPConfig  `mapstructure:"http"`
	TCP   TCPConfig   `mapstructure:"tcp"`
	V2ray V2rayConfig `mapstructure:"v2ray"`
}

// HTTPConfig HTTP服务配置
type HTTPConfig struct {
	Addr         string `mapstructure:"addr"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

// TCPConfig TCP服务配置
type TCPConfig struct {
	Addr      string `mapstructure:"addr"`
	Timeout   int    `mapstructure:"timeout"`
	KeepAlive int    `mapstructure:"keep_alive"`
}

type V2rayConfig struct {
	Addr    string `mapstructure:"addr"`
	LogPath string `mapstructure:"log_path"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type            string `mapstructure:"type"`
	DNS             string `mapstructure:"dns"`
	LogLevel        string `mapstructure:"log_level"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
	MaxRetries   int    `mapstructure:"max_retries"`
	DialTimeout  int    `mapstructure:"dial_timeout"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

// TokenConfig 令牌配置
type TokenConfig struct {
	Secret            string `mapstructure:"secret"`
	ExpireTime        int    `mapstructure:"expire_time"`
	RefreshExpireTime int    `mapstructure:"refresh_expire_time"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// Init 初始化配置
func Init(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 设置默认值
	viper.SetDefault("token.expire_time", 7*24*3600)
	viper.SetDefault("token.refresh_expire_time", 14*24*3600)

	// 绑定环境变量
	viper.AutomaticEnv()
	viper.SetEnvPrefix("NOAH")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 解析配置
	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("解析配置失败: %w", err)
	}

	return nil
}

// Get 获取配置
func Get() *Config {
	return cfg
}
