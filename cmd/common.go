package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/globalsign/mgo"
	"github.com/nawa/cryptoexchange-wallet-info/model"
	"github.com/nawa/cryptoexchange-wallet-info/storage"
	"github.com/nawa/cryptoexchange-wallet-info/storage/mongo"

	"github.com/spf13/cobra"
)

const (
	DBTimeout = time.Second * 10

	envExchangeAPIKey    = "EXCHANGE_API_KEY"
	envExchangeAPISecret = "EXCHANGE_API_SECRET"
)

type APICommand struct {
	ExchangeType string
	APIKey       string
	APISecret    string
}

type MongoCommand struct {
	MongoURL string
}

func (c *APICommand) BindArgs(cobraCmd *cobra.Command) error {
	cobraCmd.Flags().StringVarP(&c.ExchangeType, "exchange-type", "e", string(model.ExchangeTypeBittrex), fmt.Sprintf("Exchange type: [%s] (Only Bittrex is supported now)", model.ExchangeTypeBittrex))
	cobraCmd.Flags().StringVarP(&c.APIKey, "api-key", "k", "", "API Key. Can be skipped and provided by environment variable EXCHANGE_API_KEY")
	cobraCmd.Flags().StringVarP(&c.APISecret, "api-secret", "s", "", "API Secret. Can be skipped and provided by environment variable EXCHANGE_API_SECRET")
	return nil
}

func (c *APICommand) CheckArgs() error {
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
	return mongo.NewBalanceStorage(session, true), nil
}
