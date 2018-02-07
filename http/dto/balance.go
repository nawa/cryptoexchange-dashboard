package dto

import (
	"github.com/nawa/cryptoexchange-wallet-info/model"
)

type BalanceDTO struct {
	Amount     float64 `json:"amount"`
	BTCAmount  float64 `json:"btc"`
	USDTAmount float64 `json:"usdt"`
	Time       int64   `json:"time"`
}

func NewBalanceDTO(model model.CurrencyBalance) *BalanceDTO {
	return &BalanceDTO{
		Amount:     model.Amount,
		BTCAmount:  model.BTCAmount,
		USDTAmount: model.USDTAmount,
		Time:       model.Time.Unix(),
	}
}