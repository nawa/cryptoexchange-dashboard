package http

import (
	"context"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/recover"
	"github.com/nawa/cryptoexchange-dashboard/usecase"
)

type Server struct {
	app       *iris.Application
	ctx       context.Context
	ctxCancel context.CancelFunc
	log       *logrus.Entry
}

func NewServer(balanceUsecase usecase.BalanceUsecases, orderUsecase usecase.OrderUsecases) *Server {
	app := iris.New()
	app.Use(recover.New())
	app.Use(cors.Default())

	baseHandler := NewBaseHandler()
	balanceHandler := NewBalanceHandler(balanceUsecase)
	orderHandler := NewOrderHandler(orderUsecase)

	app.Get("ping", baseHandler.Ping)

	balanceGroup := app.Party("/balance")
	balanceGroup.Get("/period/hourly/{hours}", balanceHandler.Hourly)
	balanceGroup.Get("/period/weekly", balanceHandler.Weekly)
	balanceGroup.Get("/period/monthly", balanceHandler.Monthly)
	balanceGroup.Get("/period/all", balanceHandler.All)

	balanceGroup.Get("/active", balanceHandler.ActiveCurrencies)

	app.Get("/order", orderHandler.GetActiveOrders)

	server := &Server{
		app: app,
		log: logrus.WithField("component", "HTTPServer"),
	}
	return server
}

func (server *Server) Start(ctx context.Context, address string) {
	server.log.Infof("starting HTTP server on '%s'...", address)

	server.ctx, server.ctxCancel = context.WithCancel(ctx)

	err := server.app.Run(iris.Addr(address), iris.WithoutInterruptHandler)
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
