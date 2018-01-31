package storage

import "github.com/nawa/cryptoexchange-wallet-info/model"

type Exchange interface {
	GetBalance() (*model.Balance, error)
	GetMarketInfo(market string) (*model.MarketInfo, error)
	Ping() error
}
