package storage

import "github.com/nawa/cryptoexchange-dashboard/domain"

type Exchange interface {
	GetBalance() (*domain.Balance, error)
	GetMarketInfo(market string) (*domain.MarketInfo, error)
	GetOrders() ([]domain.Order, error)
	Ping() error
}
