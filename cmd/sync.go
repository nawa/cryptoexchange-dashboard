package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/globalsign/mgo"
	"github.com/spf13/cobra"

	"github.com/nawa/cryptoexchange-wallet-info/shared/exchange"
	"github.com/nawa/cryptoexchange-wallet-info/shared/storage"
	"github.com/nawa/cryptoexchange-wallet-info/shared/storage/mongo"
	"github.com/nawa/cryptoexchange-wallet-info/sync"
)

const DBTimeout = time.Second * 10

type SyncCommand struct {
	cobra.Command
	APICommand
	MongoURL   string
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
	syncCmd.APICommand.BindArgs(&syncCmd.Command)
	syncCmd.Command.Flags().StringVarP(&syncCmd.MongoURL, "db-url", "u", "", "Url to MongoDB")
	syncCmd.Command.Flags().IntVarP(&syncCmd.SyncPeriod, "period", "p", 10, "Synchronization period in sec")

	err := syncCmd.MarkFlagRequired("db-url")
	if err != nil {
		panic(err)
	}

	syncCmd.RunE = syncCmd.run
	rootCmd.AddCommand(&syncCmd.Command)
}

func (c *SyncCommand) run(_ *cobra.Command, _ []string) error {
	err := c.CheckArgs()
	if err != nil {
		return err
	}

	exchange := exchange.NewBittrexExchange(c.APIKey, c.APISecret)
	err = exchange.Ping()
	if err != nil {
		return fmt.Errorf("exchange error: %s", err)
	}

	balanceStorage, err := c.createBalanceStorage()
	if err != nil {
		return err
	}

	service := sync.NewSyncService(exchange, balanceStorage)
	ticker := sync.NewSyncTicker(time.Second*time.Duration(c.SyncPeriod), service)
	err = ticker.Start()
	if err != nil {
		return err
	}
	defer ticker.Stop()

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

func (c *SyncCommand) createBalanceStorage() (storage.BalanceStorage, error) {
	dialInfo, err := mgo.ParseURL(c.MongoURL)
	dialInfo.Timeout = DBTimeout
	if err != nil {
		return nil, fmt.Errorf("mongo URL is incorrect: %s", err)
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, fmt.Errorf("can't connect to mongo: %s", err)
	}

	return mongo.NewBalanceStorage(session, true), nil
}

func (c *SyncCommand) CheckArgs() error {
	return c.APICommand.CheckArgs()
}
