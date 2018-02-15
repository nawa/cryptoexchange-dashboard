package storage

import (
	"time"

	"github.com/nawa/cryptoexchange-dashboard/model"
)

type BalanceStorage interface {
	// Init inits the storage, such as prepares indexes and another
	Init() error
	Save(balances ...Balance) error
	FetchHourly(currency string, hours int) ([]Balance, error)
	FetchWeekly(currency string) ([]Balance, error)
	FetchMonthly(currency string) ([]Balance, error)
	FetchAll(currency string) ([]Balance, error)
	GetActiveCurrencies() ([]Balance, error)
}

type Balance struct {
	Exchange   string    `bson:"exchange"`
	Currency   string    `bson:"currency"`
	Amount     float64   `bson:"amount"`
	BTCAmount  float64   `bson:"btc_amount"`
	USDTAmount float64   `bson:"usdt_amount"`
	BTCRate    float64   `bson:"btc_rate"`
	Time       time.Time `bson:"time"`
}

func NewBalances(b *model.Balance) (result []Balance) {
	result = append(result, Balance{
		Exchange:   string(b.Exchange),
		Currency:   "total",
		Amount:     b.BTCAmount,
		BTCRate:    1,
		BTCAmount:  b.BTCAmount,
		USDTAmount: b.USDTAmount,
		Time:       b.Time,
	})

	for _, c := range b.Currencies {
		result = append(result, Balance{
			Exchange:   string(b.Exchange),
			Currency:   c.Currency,
			Amount:     c.Amount,
			BTCRate:    c.BTCRate,
			BTCAmount:  c.BTCAmount,
			USDTAmount: c.USDTAmount,
			Time:       b.Time,
		})
	}
	return result
}

func (b *Balance) ToModel() *model.CurrencyBalance {
	return &model.CurrencyBalance{
		Currency:   b.Currency,
		Amount:     b.Amount,
		BTCAmount:  b.BTCAmount,
		BTCRate:    b.BTCRate,
		USDTAmount: b.USDTAmount,
		Time:       b.Time,
	}
}
