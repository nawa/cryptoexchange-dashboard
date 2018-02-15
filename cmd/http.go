package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/nawa/cryptoexchange-wallet-info/usecase"

	log "github.com/Sirupsen/logrus"

	"github.com/nawa/cryptoexchange-wallet-info/http"
	"github.com/spf13/cobra"
)

type HTTPCommand struct {
	cobra.Command
	ExchangeAPICommand
	MongoCommand
	HTTPAddress string
}

var (
	httpCmd = &HTTPCommand{
		Command: cobra.Command{
			Use:   "http",
			Short: "REST interface to your collected wallet info",
			Long:  "Exposes REST interface to your collected wallet info",
		},
	}
)

func init() {
	err := httpCmd.MongoCommand.BindArgs(&httpCmd.Command)
	if err != nil {
		panic(err)
	}

	err = httpCmd.ExchangeAPICommand.BindArgs(&httpCmd.Command)
	if err != nil {
		panic(err)
	}

	httpCmd.Flags().StringVarP(&httpCmd.HTTPAddress, "addr", "a", "localhost:8080", "Service address")

	httpCmd.PreRunE = httpCmd.preRun
	httpCmd.RunE = httpCmd.run
	rootCmd.AddCommand(&httpCmd.Command)
}

func (c *HTTPCommand) preRun(_ *cobra.Command, _ []string) error {
	return c.ExchangeAPICommand.CheckArgs()
}

func (c *HTTPCommand) run(_ *cobra.Command, _ []string) error {
	exchange, err := c.CreateExchange()
	if err != nil {
		return err
	}

	balanceStorage, err := c.CreateBalanceStorage()
	if err != nil {
		return err
	}

	ctx, ctxCancel := context.WithCancel(context.Background())

	balanceUsecase := usecase.NewBalanceUsecase(exchange, balanceStorage)
	orderUsecase := usecase.NewOrderUsecase(exchange)

	server := http.NewServer(ctx, httpCmd.HTTPAddress, balanceUsecase, orderUsecase)

	go func() {
		defer ctxCancel()
		server.Start()
	}()

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	select {
	case <-sigC:
		ctxCancel()
	case <-ctx.Done():
	}

	server.Stop()

	log.Info("Service stopped")
	return nil
}
