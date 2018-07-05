package storage

import "github.com/nawa/cryptoexchange-dashboard/domain"

type BalanceStorage interface {
	// Init initializes the storage, such as prepares indexes and another
	Init() error
	Save(balance ...domain.Balance) error
	FetchHourly(currency string, hours int) ([]domain.Balance, error)
	FetchWeekly(currency string) ([]domain.Balance, error)
	FetchMonthly(currency string) ([]domain.Balance, error)
	FetchAll(currency string) ([]domain.Balance, error)
	GetActiveCurrencies() ([]domain.Balance, error)
}
