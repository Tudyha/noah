package server

import (
	"fmt"
	"noah/internal/server/controller"
	"noah/internal/server/dao"
	"noah/internal/server/environment"
	"noah/internal/server/gateway"
	"noah/internal/server/middleware"
	"noah/internal/server/middleware/log"
	"noah/internal/server/router"
	"noah/internal/server/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"github.com/samber/do/v2"
)

type Server struct {
	Gin *gin.Engine
	Env *environment.Environment
}

func NewServer(env *environment.Environment) *Server {
	// 依赖注入
	i := Inject()
	do.ProvideValue(i, env)

	gate, err := gateway.NewGateway(i)
	if err != nil {
		panic(err)
	}
	do.ProvideValue(i, &gate)

	gin.SetMode(gin.ReleaseMode)

	g := gin.New()
	//logger
	g.Use(log.Logger(log.SetupLogger()))

	//init routes
	router.Init(g, i)

	// panic recovery
	g.Use(gin.Recovery())

	//时间统一配置
	carbon.SetDefault(carbon.Default{
		Layout:       carbon.DateTimeLayout,
		Timezone:     carbon.PRC,
		WeekStartsAt: carbon.Sunday,
		Locale:       "zh-CN",
	})

	return &Server{
		Gin: g,
		Env: env,
	}
}

func Inject() do.Injector {
	i := do.New()
	//init db
	dao.Init(i)

	//init service
	service.Init(i)

	//init controller
	controller.Init(i)

	//init middleware
	middleware.Init(i)
	return i
}

func (s *Server) Run() {
	addr := fmt.Sprintf(":%d", s.Env.Server.Port)
	err := s.Gin.Run(addr)
	if err != nil {
		panic(err)
	}
}
