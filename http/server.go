package http

import (
	"fmt"

	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/recover"
)

type Server struct {
	port int
	iris *iris.Application
}

func NewServer(port int) *Server {
	router := iris.New()
	router.Use(recover.New())

	router.Get("ping", ping)

	server := &Server{
		port: port,
		iris: router,
	}
	return server
}

func ping(ctx iris.Context) {
	_, err := ctx.WriteString("pong")
	if err != nil {
		panic(err)
	}
}

func (server *Server) Start() error {
	return server.iris.Run(iris.Addr(fmt.Sprintf(":%d", server.port)))
}

func (server *Server) Stop() {

}
