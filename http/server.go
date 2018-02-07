package http

import (
	"context"
	"time"

	"github.com/nawa/cryptoexchange-wallet-info/storage"
	"github.com/nawa/cryptoexchange-wallet-info/usecase"

	"github.com/Sirupsen/logrus"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/recover"

	"github.com/nawa/cryptoexchange-wallet-info/http/handler"
)

type Server struct {
	addr      string
	app       *iris.Application
	ctx       context.Context
	ctxCancel context.CancelFunc
	log       *logrus.Entry
}

func NewServer(ctx context.Context, address string, balanceStorage storage.BalanceStorage) *Server {
	app := iris.New()
	app.Use(recover.New())

	baseHandler := handler.NewBaseHandler()
	balanceUsecase := usecase.NewBalanceUsecase(nil, balanceStorage)
	balanceHandler := handler.NewBalanceHandler(balanceUsecase)
	app.Get("ping", baseHandler.Ping)

	balanceGroup := app.Party("/balance")
	balanceGroup.Get("/daily", balanceHandler.Daily)

	server := &Server{
		addr: address,
		app:  app,
		log:  logrus.WithField("component", "HTTPServer"),
	}
	server.ctx, server.ctxCancel = context.WithCancel(ctx)
	return server
}

func (server *Server) Start() {
	server.log.Infof("starting HTTP server on '%s'...", server.addr)

	err := server.app.Run(iris.Addr(server.addr), iris.WithoutInterruptHandler)
	if err != nil {
		server.log.WithError(err).Errorf("HTTP server interrupted with error")
	}
}

func (server *Server) Stop() {
	ctx, cancel := context.WithTimeout(server.ctx, time.Second*5)
	defer cancel()
	err := server.app.Shutdown(ctx)
	if err != nil {
		server.log.WithError(err).Errorf("HTTP server stopped with error")
	} else {
		server.log.Info("HTTP server stopped")
	}
}
