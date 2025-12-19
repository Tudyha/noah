package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"noah/internal/controller"
	"noah/internal/dao"
	"noah/internal/database"
	"noah/internal/server"
	"noah/internal/service"
	"noah/internal/session"
	"noah/pkg/app"
	"noah/pkg/config"
	"noah/pkg/logger"
)

// @title Noah API
// @version 1.0
// @description This is a sample server for Noah API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// 初始化配置
	if err := config.Init("config/config.yaml"); err != nil {
		panic(fmt.Sprintf("初始化配置失败: %v", err))
	}

	cfg := config.Get()

	// 初始化日志
	if err := logger.Init(&cfg.Logger); err != nil {
		panic(fmt.Sprintf("初始化日志失败: %v", err))
	}
	defer logger.Sync()

	// 设置Gin模式
	if cfg.Server.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化数据库
	dbConfig := cfg.Database
	dbConfig.LogLevel = cfg.Logger.Level
	if err := database.Init(&dbConfig); err != nil {
		logger.Error("初始化数据库失败", "err", err)
		os.Exit(1)
	}
	defer database.CloseDB()

	// 初始化dao层
	if err := dao.Init(); err != nil {
		logger.Error("初始化dao层失败", "err", err)
		os.Exit(1)
	}

	// 初始化service层
	if err := service.Init(); err != nil {
		logger.Error("初始化service层失败", "err", err)
		os.Exit(1)
	}

	// 初始化controller层
	if err := controller.Init(); err != nil {
		logger.Error("初始化controller层失败", "err", err)
		os.Exit(1)
	}

	// 初始化tcp连接管理器
	if err := session.Init(); err != nil {
		logger.Error("初始化tcp连接管理器失败", "err", err)
		os.Exit(1)
	}

	// 创建服务实例
	servers := []app.Server{
		server.NewHTTPServer(),
		server.NewTCPServer(),
	}

	// 启动应用
	app := app.NewApp(servers...)
	app.Run()
}
