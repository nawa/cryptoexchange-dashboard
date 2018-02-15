package http

import (
	"context"
	"time"

	"github.com/nawa/cryptoexchange-wallet-info/http/handler"
	"github.com/nawa/cryptoexchange-wallet-info/usecase"

	"github.com/Sirupsen/logrus"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/recover"
)

type Server struct {
	addr      string
	app       *iris.Application
	ctx       context.Context
	ctxCancel context.CancelFunc
	log       *logrus.Entry
}

func NewServer(ctx context.Context, address string, balanceUsecase usecase.BalanceUsecase, orderUsecase usecase.OrderUsecase) *Server {
	app := iris.New()
	app.Use(recover.New())

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})
	app.Use(crs)

	baseHandler := handler.NewBaseHandler()
	balanceHandler := handler.NewBalanceHandler(balanceUsecase)
	orderHandler := handler.NewOrderHandler(orderUsecase)

	app.Get("ping", baseHandler.Ping)

	balanceGroup := app.Party("/balance")
	balanceGroup.Get("/period/hourly/{hours:int}", balanceHandler.Hourly)
	balanceGroup.Get("/period/weekly", balanceHandler.Weekly)
	balanceGroup.Get("/period/monthly", balanceHandler.Monthly)
	balanceGroup.Get("/period/all", balanceHandler.All)

	balanceGroup.Get("/active", balanceHandler.ActiveCurrencies)

	app.Get("/order", orderHandler.GetActiveOrders)

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
