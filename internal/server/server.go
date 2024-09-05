package server

import (
	"github.com/golang-module/carbon/v2"
	"noah/internal/server/dao"
	"noah/internal/server/environment"
	"noah/internal/server/middleware"
	"noah/internal/server/middleware/log"
	"noah/internal/server/routes"
	"noah/internal/server/service"

	"github.com/gin-gonic/gin"
)

type Server struct {
	G   *gin.Engine
	Env environment.Environment
}

func NewServer() *Server {
	//加载配置文件
	env, err := environment.LoadEnvironment()
	if err != nil {
		panic(err)
	}

	middleware.SetAdminPassword(env.Admin.Password)

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

	//cron
	middleware.LoadCron()

	//load service
	service.LoadService()

	gin.SetMode(gin.ReleaseMode)
	g := gin.New()
	//log
	logger := log.SetupLogger()
	g.Use(log.Logger(logger))

	//init routes
	r := routes.NewRouter(g)
	r.LoadRoutes()

	// panic recovery
	g.Use(gin.Recovery())

	return &Server{
		G:   g,
		Env: env,
	}
}

func (s *Server) Run() {
	err := s.G.Run(":" + s.Env.Server.Port)
	if err != nil {
		panic(err)
	}
}
