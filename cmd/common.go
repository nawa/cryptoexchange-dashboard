package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

const (
	envExchangeAPIKey    = "EXCHANGE_API_KEY"
	envExchangeAPISecret = "EXCHANGE_API_SECRET"
)

type APICommand struct {
	ExchangeType string
	APIKey       string
	APISecret    string
}

func (c *APICommand) BindArgs(cobraCmd *cobra.Command) {
	cobraCmd.Flags().StringVarP(&c.ExchangeType, "exchange-type", "e", "bittrex", "Exchange type: [bittrex] (Only Bittrex is supported now)")
	cobraCmd.Flags().StringVarP(&c.APIKey, "api-key", "k", "", "API Key. Can be skipped and provided by environment variable EXCHANGE_API_KEY")
	cobraCmd.Flags().StringVarP(&c.APISecret, "api-secret", "s", "", "API Secret. Can be skipped and provided by environment variable EXCHANGE_API_SECRET")
}

func (c *APICommand) CheckArgs() error {
	if c.ExchangeType != "bittrex" {
		return errors.New("--exchange-type is wrong, supported values: [bittrex] (Only Bittrex is supported now)")
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
