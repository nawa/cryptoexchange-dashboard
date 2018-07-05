package testdata

import (
	"time"

	"github.com/nawa/cryptoexchange-dashboard/domain"
)

func Balances() []domain.Balance {
	return []domain.Balance{
		{
			Exchange:   domain.ExchangeTypeBittrex,
			Currency:   "CUR1",
			Amount:     100,
			BTCAmount:  200,
			USDTAmount: 300,
			Time:       time.Unix(0, 0).UTC(),
		}, {
			Exchange:   domain.ExchangeTypeBittrex,
			Currency:   "CUR2",
			Amount:     400,
			BTCAmount:  500,
			USDTAmount: 600,
			Time:       time.Unix(0, 0).UTC().Add(time.Hour),
		},
	}
}

func BalancesWithTotal() []domain.Balance {
	return append(Balances(), domain.Balance{
		Exchange:   domain.ExchangeTypeBittrex,
		Currency:   "total",
		Amount:     0,
		BTCAmount:  700,
		USDTAmount: 900,
		Time:       Balances()[0].Time,
	})
}

func Orders() []domain.Order {
	return []domain.Order{
		{
			Exchange:    domain.ExchangeTypeBittrex,
			Market:      "market",
			SellNowRate: 0.111,
			Time:        time.Unix(0, 0).UTC(),
			USDTRate:    0.999,
		},
		{
			Exchange:    domain.ExchangeTypeBittrex,
			Market:      "market2",
			SellNowRate: 0.22222,
			Time:        time.Unix(0, 0).UTC().Add(time.Hour),
			USDTRate:    0.88888,
		},
	}
}
