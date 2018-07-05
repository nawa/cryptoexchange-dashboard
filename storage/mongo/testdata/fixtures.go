package testdata

import (
	"github.com/nawa/cryptoexchange-dashboard/domain"
)

func Balances() []domain.Balance {
	return []domain.Balance{
		{
			Exchange:   domain.ExchangeTypeBittrex,
			Currency:   "CUR1",
			Amount:     1,
			BTCAmount:  2,
			USDTAmount: 3,
		},
		{
			Exchange:   domain.ExchangeTypeBittrex,
			Currency:   "CUR2",
			Amount:     4,
			BTCAmount:  5,
			USDTAmount: 6,
		},
		{
			Exchange:   domain.ExchangeTypeBittrex,
			Currency:   "total",
			BTCAmount:  100,
			USDTAmount: 1000,
		},
		{
			Exchange:   domain.ExchangeTypeBittrex,
			Currency:   "CUR1",
			Amount:     1,
			BTCAmount:  2,
			USDTAmount: 3,
		},
		{
			Exchange:   domain.ExchangeTypeBittrex,
			Currency:   "CUR2",
			Amount:     4,
			BTCAmount:  5,
			USDTAmount: 6,
		},
		{
			Exchange:   domain.ExchangeTypeBittrex,
			Currency:   "total",
			BTCAmount:  100,
			USDTAmount: 1000,
		},
	}

}
