package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0xAX/notificator"
	log "github.com/Sirupsen/logrus"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"

	"github.com/nawa/cryptoexchange-dashboard/storage"
	"github.com/nawa/cryptoexchange-dashboard/storage/exchange"
)

type NotifyCommand struct {
	cobra.Command
	ExchangeAPICommand
	Market        string
	GreaterThan   float64
	LessThan      float64
	RefreshPeriod int
}

var (
	notifyCmd = &NotifyCommand{
		Command: cobra.Command{
			Use:   "notify",
			Short: "Notifies when price of coin is reached some value",
			Long:  "Notifies when price of coin is reached some value. \nATTENTION: would be more secure is to generate keys with readonly permission",
		},
	}
)

func init() {
	err := notifyCmd.ExchangeAPICommand.BindArgs(&notifyCmd.Command)
	if err != nil {
		panic(err)
	}
	notifyCmd.Command.Flags().IntVarP(&notifyCmd.RefreshPeriod, "period", "p", 10, "Refresh period in sec")
	notifyCmd.Command.Flags().StringVarP(&notifyCmd.Market, "market", "m", "", "Market name, for example 'BTC-ETH'")
	notifyCmd.Command.Flags().Float64Var(&notifyCmd.GreaterThan, "gt", 0, "Notify when price is greater than value")
	notifyCmd.Command.Flags().Float64Var(&notifyCmd.LessThan, "lt", 0, "Notify when price is less than value")

	err = notifyCmd.MarkFlagRequired("market")
	if err != nil {
		panic(err)
	}

	notifyCmd.PreRunE = notifyCmd.preRun
	notifyCmd.RunE = notifyCmd.run
	rootCmd.AddCommand(&notifyCmd.Command)
}

func (c *NotifyCommand) preRun(_ *cobra.Command, _ []string) error {
	err := c.ExchangeAPICommand.CheckArgs()
	if err != nil {
		return err
	}

	if c.GreaterThan == 0 && c.LessThan == 0 {
		return errors.New("--gt or --lt must be defined")
	}

	if c.GreaterThan > 0 && c.LessThan > 0 {
		return errors.New("only one of --gt or --lt must be defined")
	}

	if c.GreaterThan < 0 {
		return errors.New("--gt must be (0, ∞)")
	} else if c.LessThan < 0 {
		return errors.New("--lt must be (0, ∞)")
	}
	return nil
}

func (c *NotifyCommand) run(_ *cobra.Command, _ []string) error {
	exchange := exchange.NewBittrexExchange(c.APIKey, c.APISecret)
	err := exchange.Ping()
	if err != nil {
		return err
	}
	ticker := time.NewTicker(time.Second * time.Duration(c.RefreshPeriod))

	resultCh := make(chan error, 1)
	go func() {
		for range ticker.C {
			lastPrice, err := c.checkMarketLastPrice(exchange)
			if err != nil {
				log.Error(err)
				continue
			}

			if lastPrice != nil {
				ticker.Stop()
				err = c.sendNotification(*lastPrice)
				resultCh <- err
				return
			}
		}
	}()

	exitC := make(chan os.Signal, 1)
	signal.Notify(exitC,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	select {
	case <-exitC:
		fmt.Println("Shutting down...")
	case err := <-resultCh:
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (c *NotifyCommand) checkMarketLastPrice(exchange storage.Exchange) (*float64, error) {
	marketInfo, err := exchange.GetMarketInfo(c.Market)
	if err != nil {
		return nil, err
	}

	if (notifyCmd.GreaterThan > 0 && marketInfo.Last >= notifyCmd.GreaterThan) ||
		(notifyCmd.LessThan > 0 && marketInfo.Last <= notifyCmd.LessThan) {
		return &marketInfo.Last, nil
	}
	return nil, nil
}

func (c *NotifyCommand) sendNotification(lastPrice float64) error {
	notifier := notificator.New(notificator.Options{
		DefaultIcon: "",
		AppName:     "cryptoexchange-dashboard",
	})

	msg := fmt.Sprintf("%s is reached price %s", notifyCmd.Market, decimal.NewFromFloat(lastPrice))
	return notifier.Push("Coin price notifier", msg, "", notificator.UR_CRITICAL)
}
