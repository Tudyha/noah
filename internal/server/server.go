package server

import (
	"noah/internal/server/dao"
	"noah/internal/server/environment"
	"noah/internal/server/middleware"
	"noah/internal/server/routes"

	"github.com/gin-gonic/gin"
	_ "noah/internal/server/logic/client"
	_ "noah/internal/server/logic/pty"
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

	g := gin.Default()
	// panic recovery
	//g.Use(gin.Recovery())

	//init routes
	r := routes.NewRouter(g)
	r.LoadRoutes()

	//cron
	middleware.LoadCron()

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
