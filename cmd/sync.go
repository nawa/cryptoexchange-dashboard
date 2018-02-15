package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nawa/cryptoexchange-dashboard/usecase"

	"github.com/spf13/cobra"
)

type SyncCommand struct {
	cobra.Command
	ExchangeAPICommand
	MongoCommand
	SyncPeriod int
}

var (
	syncCmd = &SyncCommand{
		Command: cobra.Command{
			Use:   "sync",
			Short: "Syncs your exchange data",
			Long:  "Syncs your exchange data with database in background. \nATTENTION: would be more secure is to generate keys with readonly permission",
		},
	}
)

func init() {
	err := syncCmd.ExchangeAPICommand.BindArgs(&syncCmd.Command)
	if err != nil {
		panic(err)
	}
	err = syncCmd.MongoCommand.BindArgs(&syncCmd.Command)
	if err != nil {
		panic(err)
	}
	syncCmd.Command.Flags().IntVarP(&syncCmd.SyncPeriod, "period", "p", 10, "Synchronization period in sec")

	syncCmd.PreRunE = syncCmd.preRun
	syncCmd.RunE = syncCmd.run
	rootCmd.AddCommand(&syncCmd.Command)
}

func (c *SyncCommand) preRun(_ *cobra.Command, _ []string) error {
	return c.ExchangeAPICommand.CheckArgs()
}

func (c *SyncCommand) run(_ *cobra.Command, _ []string) error {
	exchange, err := c.CreateExchange()
	if err != nil {
		return err
	}

	balanceStorage, err := c.CreateBalanceStorage()
	if err != nil {
		return err
	}

	balanceUsecase := usecase.NewBalanceUsecase(exchange, balanceStorage)
	synchronizer, err := balanceUsecase.StartSyncFromExchangePeriodically(c.SyncPeriod)
	if err != nil {
		return err
	}
	defer synchronizer.Stop()

	exitC := make(chan os.Signal, 1)
	signal.Notify(exitC,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-exitC
	fmt.Println("Shutting down...")

	return nil
}
