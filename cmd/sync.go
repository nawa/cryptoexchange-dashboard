package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/nawa/cryptoexchange-wallet-info/storage/exchange"
	"github.com/nawa/cryptoexchange-wallet-info/usecase"
)

type SyncCommand struct {
	cobra.Command
	APICommand
	MongoCommand
	SyncPeriod int
}

var (
	syncCmd = &SyncCommand{
		Command: cobra.Command{
			Use:   "sync",
			Short: "Syncs your wallet data",
			Long:  "Syncs your wallet data with database in background. \nATTENTION: would be more secure is to generate keys with readonly permission",
		},
	}
)

func init() {
	err := syncCmd.APICommand.BindArgs(&syncCmd.Command)
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
	return c.APICommand.CheckArgs()
}

func (c *SyncCommand) run(_ *cobra.Command, _ []string) error {
	exchange := exchange.NewBittrexExchange(c.APIKey, c.APISecret)
	err := exchange.Ping()
	if err != nil {
		return fmt.Errorf("exchange error: %s", err)
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
