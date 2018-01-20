package storage

import "time"

type BalanceStorage interface {
	Save(balances ...Balance) error
	Find() ([]Balance, error)
}

type Balance struct {
	Exchange   string    `bson:"exchange"`
	Currency   string    `bson:"currency"`
	Amount     float64   `bson:"amount"`
	BTCAmount  float64   `bson:"btc_amount"`
	BTCRate    float64   `bson:"btc_rate"`
	USDTAmount float64   `bson:"usdt_rate"`
	Time       time.Time `bson:"time"`
}
