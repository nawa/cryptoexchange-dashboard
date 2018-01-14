package model

import (
	"time"

	"github.com/nawa/cryptoexchange-wallet-info/shared/storage"
)

type CurrencyBalance struct {
	Currency   string
	Amount     float64
	BTCAmount  float64
	BTCRate    float64
	USDTAmount float64
}

type Balance struct {
	Currencies []CurrencyBalance
	BTCAmount  float64
	USDTAmount float64
	Time       time.Time
}

func (b *Balance) ToStorageModel() (result []storage.Balance) {
	result = append(result, storage.Balance{
		Currency:   "total",
		Amount:     b.BTCAmount,
		BTCRate:    1,
		BTCAmount:  b.BTCAmount,
		USDTAmount: b.USDTAmount,
		Time:       b.Time,
	})

	for _, c := range b.Currencies {
		result = append(result, storage.Balance{
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
