package testdata

import (
	"time"

	"github.com/nawa/cryptoexchange-dashboard/model"
	"github.com/nawa/cryptoexchange-dashboard/storage"
)

func CurrencyBalances() []model.CurrencyBalance {
	return []model.CurrencyBalance{
		{
			Amount:     1,
			BTCAmount:  1,
			Currency:   "total",
			Time:       time.Unix(0, 0).UTC(),
			USDTAmount: 2,
		},
		{
			Amount:     10,
			BTCAmount:  9,
			Currency:   "CUR",
			Time:       time.Unix(0, 0).UTC().Add(time.Hour),
			USDTAmount: 7,
		},
	}
}

func ModelBalance() *model.Balance {
	return &model.Balance{
		Currencies: []model.CurrencyBalance{CurrencyBalances()[1]},
		Exchange:   model.ExchangeTypeBittrex,
		BTCAmount:  1,
		USDTAmount: 2,
		Time:       time.Unix(0, 0).UTC(),
	}
}

func StorageBalances() []storage.Balance {
	modelBalance := ModelBalance()
	return []storage.Balance{
		{
			Amount:     modelBalance.BTCAmount,
			BTCAmount:  modelBalance.BTCAmount,
			Currency:   "total",
			Exchange:   string(modelBalance.Exchange),
			Time:       modelBalance.Time,
			USDTAmount: modelBalance.USDTAmount,
		}, {
			Amount:     modelBalance.Currencies[0].Amount,
			BTCAmount:  modelBalance.Currencies[0].BTCAmount,
			Currency:   modelBalance.Currencies[0].Currency,
			Exchange:   string(modelBalance.Exchange),
			Time:       modelBalance.Currencies[0].Time,
			USDTAmount: modelBalance.Currencies[0].USDTAmount,
		}}
}

func Orders() []model.Order {
	return []model.Order{
		{
			Exchange:    model.ExchangeTypeBittrex,
			Market:      "market",
			SellNowRate: 0.111,
			Time:        time.Unix(0, 0).UTC(),
			USDTRate:    0.999,
		},
		{
			Exchange:    model.ExchangeTypeBittrex,
			Market:      "market2",
			SellNowRate: 0.22222,
			Time:        time.Unix(0, 0).UTC().Add(time.Hour),
			USDTRate:    0.88888,
		},
	}
}
