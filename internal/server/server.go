package server

import (
	"noah/internal/server/dao"
	"noah/internal/server/environment"
	"noah/internal/server/gateway"
	"noah/internal/server/middleware"
	"noah/internal/server/middleware/log"
	"noah/internal/server/routes"
	"noah/internal/server/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"github.com/samber/do/v2"
)

type Server struct {
	Gin *gin.Engine
	Env environment.Environment
}

func NewServer() *Server {
	//加载配置文件
	env, err := environment.LoadEnvironment()
	if err != nil {
		panic(err)
	}
	// 依赖注入
	injector := do.New()
	Inject(injector)
	do.ProvideValue(injector, env)

	gin.SetMode(gin.ReleaseMode)
	g := gin.New()
	//logger
	logger := log.SetupLogger()
	g.Use(log.Logger(logger))

	//init routes
	r := routes.NewRouter(g, injector)
	r.LoadRoutes()

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

func Inject(i do.Injector) {
	//init db
	dbError := dao.Init(i)
	if dbError != nil {
		panic(dbError)
	}
	//gateway
	do.Provide(i, gateway.NewGateway)
	//load service
	service.LoadService(i)
	//cron
	middleware.LoadCron()
}

func (s *Server) Run() {
	addr := ":" + s.Env.Server.Port
	err := s.Gin.Run(addr)
	if err != nil {
		panic(err)
	}
}
