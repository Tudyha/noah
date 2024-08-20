package server

import (
	"noah/internal/server/dao"
	"noah/internal/server/environment"
	"noah/internal/server/routes"

	"github.com/gin-gonic/gin"
	_ "noah/internal/server/logic/client"
	_ "noah/internal/server/logic/device"
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

	//init db
	dbError := dao.InitDb(env.Database)
	if dbError != nil {
		panic(dbError)
	}

	g := gin.Default()

	//init routes
	r := routes.NewRouter(g)
	r.LoadRoutes()

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
