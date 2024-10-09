package server

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"github.com/samber/do/v2"
	"noah/internal/server/dao"
	"noah/internal/server/environment"
	"noah/internal/server/gateway"
	"noah/internal/server/middleware"
	"noah/internal/server/middleware/log"
	"noah/internal/server/routes"
	"noah/internal/server/service"
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

	middleware.SetAdminPassword(env.Admin.Password)

	// 依赖注入
	injector := do.New()
	Inject(injector)

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

	//init db
	dbError := dao.InitDb(env.Database)
	if dbError != nil {
		panic(dbError)
	}

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
