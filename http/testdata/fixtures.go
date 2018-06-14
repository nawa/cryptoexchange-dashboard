package testdata

import (
	"time"

	"github.com/nawa/cryptoexchange-dashboard/domain"
)

func CurrencyBalances() map[string][]domain.CurrencyBalance {
	return map[string][]domain.CurrencyBalance{
		"CUR1": {
			{
				Currency:   "CUR1",
				Amount:     1,
				BTCAmount:  1,
				Time:       time.Unix(0, 0).UTC(),
				USDTAmount: 2,
			},
			{
				Currency:   "CUR1",
				Amount:     2,
				BTCAmount:  2,
				Time:       time.Unix(0, 0).UTC().Add(time.Hour),
				USDTAmount: 4,
			},
		},
		"CUR2": {
			{
				Currency:   "CUR2",
				Amount:     3,
				BTCAmount:  3,
				Time:       time.Unix(0, 0).UTC().Add(2 * time.Hour),
				USDTAmount: 5,
			},
		},
	}
}

func Orders() []domain.Order {
	return []domain.Order{
		{
			Exchange:    domain.ExchangeTypeBittrex,
			Market:      "market1",
			Time:        time.Unix(0, 0).UTC(),
			BuyRate:     1,
			Amount:      2,
			SellNowRate: 3,
			USDTRate:    4,
		},
		{
			Exchange:    domain.ExchangeTypeBittrex,
			Market:      "market2",
			Time:        time.Unix(0, 0).UTC().Add(time.Hour),
			BuyRate:     5,
			Amount:      6,
			SellNowRate: 7,
			USDTRate:    8,
		},
	}
}
