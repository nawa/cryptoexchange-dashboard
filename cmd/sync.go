package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nawa/cryptoexchange-wallet-info/sync"
	"github.com/nawa/cryptoexchange-wallet-info/sync/exchange"
	"github.com/spf13/cobra"
)

const (
	envBittrexAPIKey    = "BITTREX_API_KEY"
	envBittrexAPISecret = "BITTREX_API_SECRET"
)

var (
	mongoURL     string
	exchangeType string
	apiKey       string
	apiSecret    string
	period       int

	syncCmd = &cobra.Command{
		Use:   "sync",
		Short: "Syncs your wallet data",
		Long:  "Syncs your wallet data with database in background. \nATTENTION: would be more secure is to generate keys with readonly permission",
		RunE:  syncCmdRun,
	}
)

func init() {

	syncCmd.Flags().StringVarP(&mongoURL, "db-url", "u", "bittrex", "Url to MongoDB")
	syncCmd.Flags().StringVarP(&exchangeType, "exchange-type", "e", "", "Exchange type: [bittrex] (Only Bittrex is supported now)")
	syncCmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "API Key. Can be skipped and provided by environment variable BITTREX_API_KEY")
	syncCmd.Flags().StringVarP(&apiSecret, "api-secret", "s", "", "API Secret. Can be skipped and provided by environment variable BITTREX_API_SECRET")
	syncCmd.Flags().IntVarP(&period, "period", "p", 10, "Synchronization period in sec")

	err := syncCmd.MarkFlagRequired("db-url")
	if err != nil {
		panic(err)
	}

	rootCmd.AddCommand(syncCmd)
}

func syncCmdRun(_ *cobra.Command, _ []string) error {
	err := checkRequiredArgs()
	if err != nil {
		return err
	}

	exchange := exchange.NewBittrexExchange(apiKey, apiSecret)
	err = exchange.Ping()
	if err != nil {
		return fmt.Errorf("exchange error: %s", err)
	}

	service := sync.NewSyncService(exchange)
	ticker := sync.NewSyncTicker(time.Second*time.Duration(period), service)
	err = ticker.Start()
	if err != nil {
		return err
	}
	defer ticker.Stop()

	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-c
	fmt.Println("Shutting down...")

	return nil
}

func checkRequiredArgs() error {
	if apiKey == "" {
		apiKey = os.Getenv(envBittrexAPIKey)
		if apiKey == "" {
			return errors.New("--api-key argument or 'BITTREX_API_KEY' environment variable must be provided")
		}
	}

	if apiSecret == "" {
		apiSecret = os.Getenv(envBittrexAPISecret)
		if apiSecret == "" {
			return errors.New("--api-secret argument or 'BITTREX_API_SECRET' environment variable must be provided")
		}
	}
	return nil
}
