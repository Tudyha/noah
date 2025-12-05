package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"noah/pkg/app"
	"noah/pkg/config"

	"noah/internal/api"
	"noah/internal/middleware"
	"noah/internal/service"
	"noah/pkg/logger"
)

type httpServer struct {
	s *http.Server
}

func NewHTTPServer() app.Server {
	cfg := config.Get()

	router := gin.New()

	authService := service.GetAuthService()

	// 注册中间件
	registerMiddlewares(router, authService)

	// 注册路由
	api.RegisterRoutes(router)

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Server.HTTP.Port),
		Handler:        router,
		ReadTimeout:    time.Duration(cfg.Server.HTTP.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.Server.HTTP.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &httpServer{s: s}
}

func (h *httpServer) Start(ctx context.Context) error {
	logger.Info("http server start: ", "addr", h.s.Addr)
	if err := h.s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("HTTP服务启动失败", "err", err)
		return err
	}
	return nil
}

func (h *httpServer) Stop(ctx context.Context) error {
	if h.s == nil {
		return nil
	}
	return h.s.Shutdown(ctx)
}

func (h *httpServer) String() string {
	return "HTTP Server"
}

// registerMiddlewares 注册中间件
func registerMiddlewares(router *gin.Engine, authService service.AuthService) {
	// 基础中间件
	router.Use(gin.Recovery())

	// 鉴权中间件
	router.Use(middleware.Auth(authService))

	// TODO: 注册其他中间件
	// router.Use(middleware.Logger())
	// router.Use(middleware.CORS())
	// router.Use(middleware.RequestID())
	// router.Use(middleware.RateLimit())
}
