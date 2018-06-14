package dto

import "github.com/nawa/cryptoexchange-dashboard/domain"

type BalanceDTO map[string][]CurrencyBalanceDTO

type CurrencyBalanceDTO struct {
	Amount     float64 `json:"amount"`
	BTCAmount  float64 `json:"btc"`
	USDTAmount float64 `json:"usdt"`
	Time       int64   `json:"time"`
}

func NewCurrencyBalanceDTO(model domain.CurrencyBalance) *CurrencyBalanceDTO {
	return &CurrencyBalanceDTO{
		Amount:     model.Amount,
		BTCAmount:  model.BTCAmount,
		USDTAmount: model.USDTAmount,
		Time:       model.Time.Unix(),
	}
}

func (b BalanceDTO) Add(currency string, balance ...CurrencyBalanceDTO) {
	if len(balance) == 0 {
		return
	}
	b[currency] = append(b[currency], balance...)
}
