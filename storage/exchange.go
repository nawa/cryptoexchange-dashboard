package storage

import "github.com/nawa/cryptoexchange-dashboard/model"

type Exchange interface {
	GetBalance() (*model.Balance, error)
	GetMarketInfo(market string) (*model.MarketInfo, error)
	GetOrders() ([]model.Order, error)
	Ping() error
}
