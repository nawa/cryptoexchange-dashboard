package model

import (
	"time"
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
	Exchange   ExchangeType
	BTCAmount  float64
	USDTAmount float64
	Time       time.Time
}
