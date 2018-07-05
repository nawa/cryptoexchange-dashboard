package dto

import "github.com/nawa/cryptoexchange-dashboard/domain"

type BalancesResponse map[string][]BalanceDTO //currency/balances for time range

type BalanceDTO struct {
	Amount     float64 `json:"amount"`
	BTCAmount  float64 `json:"btc"`
	USDTAmount float64 `json:"usdt"`
	Time       int64   `json:"time"`
}

func NewBalanceDTO(model domain.Balance) *BalanceDTO {
	return &BalanceDTO{
		Amount:     model.Amount,
		BTCAmount:  model.BTCAmount,
		USDTAmount: model.USDTAmount,
		Time:       model.Time.Unix(),
	}
}

func (b BalancesResponse) Add(currency string, balance ...BalanceDTO) {
	if len(balance) == 0 {
		return
	}
	b[currency] = append(b[currency], balance...)
}
