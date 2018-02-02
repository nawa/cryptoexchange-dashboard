package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"

	"github.com/nawa/cryptoexchange-wallet-info/http"
	"github.com/spf13/cobra"
)

type HTTPCommand struct {
	cobra.Command
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

	httpCmd.Flags().StringVarP(&httpCmd.HTTPAddress, "addr", "a", "localhost:8080", "Service address")

	httpCmd.RunE = httpCmd.run
	rootCmd.AddCommand(&httpCmd.Command)
}

func (c *HTTPCommand) run(_ *cobra.Command, _ []string) error {
	ctx, ctxCancel := context.WithCancel(context.Background())

	server := http.NewServer(ctx, httpCmd.HTTPAddress)

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
