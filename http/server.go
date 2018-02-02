package http

import (
	"context"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/recover"

	"github.com/nawa/cryptoexchange-wallet-info/http/handler"
)

type Server struct {
	addr      string
	app       *iris.Application
	ctx       context.Context
	ctxCancel context.CancelFunc
}

func NewServer(ctx context.Context, address string) *Server {
	app := iris.New()
	app.Use(recover.New())

	baseHandler := handler.NewBaseHandler()
	balanceHandler := handler.NewBalanceHandler(nil)
	app.Get("ping", baseHandler.Ping)

	balanceGroup := app.Party("/balance")
	balanceGroup.Get("/", balanceHandler.Get)

	server := &Server{
		addr: address,
		app:  app,
	}
	server.ctx, server.ctxCancel = context.WithCancel(ctx)
	return server
}

func (server *Server) Start() {
	log.Infof("starting HTTP server on '%s'...", server.addr)

	err := server.app.Run(iris.Addr(server.addr), iris.WithoutInterruptHandler)
	if err != nil {
		log.Errorf("HTTP server interrupted with error: %s", err)
	}
}

func (server *Server) Stop() {
	ctx, cancel := context.WithTimeout(server.ctx, time.Second*5)
	defer cancel()
	err := server.app.Shutdown(ctx)
	if err != nil {
		log.Errorf("HTTP server stopped with error: %s", err)
	} else {
		log.Info("HTTP server stopped")
	}
}
