package storage

import "github.com/nawa/cryptoexchange-dashboard/domain"

type BalanceStorage interface {
	// Init initializes the storage, such as prepares indexes and another
	Init() error
	Save(*domain.Balance) error
	FetchHourly(currency string, hours int) ([]domain.CurrencyBalance, error)
	FetchWeekly(currency string) ([]domain.CurrencyBalance, error)
	FetchMonthly(currency string) ([]domain.CurrencyBalance, error)
	FetchAll(currency string) ([]domain.CurrencyBalance, error)
	GetActiveCurrencies() ([]domain.CurrencyBalance, error)
}
