package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/nawa/cryptoexchange-dashboard/model"
	"github.com/nawa/cryptoexchange-dashboard/storage"
	"github.com/nawa/cryptoexchange-dashboard/storage/exchange"
	"github.com/nawa/cryptoexchange-dashboard/storage/mongo"

	"github.com/Sirupsen/logrus"
	"github.com/globalsign/mgo"
	"github.com/spf13/cobra"
)

const (
	DBTimeout = time.Second * 10 //TODO Indexing takes more than 10sec - i/o timeout error

	envExchangeAPIKey    = "EXCHANGE_API_KEY"
	envExchangeAPISecret = "EXCHANGE_API_SECRET"
)

type ExchangeAPICommand struct {
	ExchangeType string
	APIKey       string
	APISecret    string
}

type MongoCommand struct {
	MongoURL string
}

func (c *ExchangeAPICommand) BindArgs(cobraCmd *cobra.Command) error {
	cobraCmd.Flags().StringVarP(&c.ExchangeType, "exchange-type", "e", string(model.ExchangeTypeBittrex), fmt.Sprintf("Exchange type: [%s] (Only Bittrex is supported now)", model.ExchangeTypeBittrex))
	cobraCmd.Flags().StringVarP(&c.APIKey, "api-key", "k", "", "API Key. Can be skipped and provided by environment variable EXCHANGE_API_KEY")
	cobraCmd.Flags().StringVarP(&c.APISecret, "api-secret", "s", "", "API Secret. Can be skipped and provided by environment variable EXCHANGE_API_SECRET")
	return nil
}

func (c *ExchangeAPICommand) CheckArgs() error {
	if c.ExchangeType != string(model.ExchangeTypeBittrex) {
		return fmt.Errorf("--exchange-type is wrong, supported values: [%s] (Only Bittrex is supported now)", model.ExchangeTypeBittrex)
	}

	if c.APIKey == "" {
		c.APIKey = os.Getenv(envExchangeAPIKey)
		if c.APIKey == "" {
			return errors.New("--api-key argument or 'EXCHANGE_API_KEY' environment variable must be provided")
		}
	}

	if c.APISecret == "" {
		c.APISecret = os.Getenv(envExchangeAPISecret)
		if c.APISecret == "" {
			return errors.New("--api-secret argument or 'EXCHANGE_API_SECRET' environment variable must be provided")
		}
	}
	return nil
}

func (c *ExchangeAPICommand) CreateExchange() (storage.Exchange, error) {
	exchange := exchange.NewBittrexExchange(c.APIKey, c.APISecret)
	err := exchange.Ping()
	if err != nil {
		return nil, fmt.Errorf("exchange error: %s", err)
	}
	return exchange, nil
}

func (c *MongoCommand) BindArgs(cobraCmd *cobra.Command) error {
	cobraCmd.Flags().StringVarP(&c.MongoURL, "db-url", "u", "", "Url to MongoDB")

	return cobraCmd.MarkFlagRequired("db-url")
}

func (c *MongoCommand) createMongoSession() (*mgo.Session, error) {
	dialInfo, err := mgo.ParseURL(c.MongoURL)
	dialInfo.Timeout = DBTimeout
	if err != nil {
		return nil, fmt.Errorf("mongo URL is incorrect: %s", err)
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, fmt.Errorf("can't connect to mongo: %s", err)
	}

	return session, nil
}

func (c *MongoCommand) CreateBalanceStorage() (storage.BalanceStorage, error) {
	session, err := c.createMongoSession()
	if err != nil {
		return nil, err
	}
	balanceStorage := mongo.NewBalanceStorage(session, true)
	go func() {
		err = balanceStorage.Init()
		if err != nil {
			logrus.WithField("component", "MongoCommand").
				WithError(err).
				Fatal("balance storage initialization error")
		}
	}()

	return balanceStorage, nil
}
