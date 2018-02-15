package dto

import (
	"github.com/nawa/cryptoexchange-wallet-info/model"
)

type OrderDTO struct {
	Market      string  `json:"market"`
	MarketLink  string  `json:"market_link"`
	Time        int64   `json:"time"`
	BuyRate     float64 `json:"buy_rate"`
	Amount      float64 `json:"amount"`
	SellNowRate float64 `json:"sellnow_rate"`
	USDTRate    float64 `json:"usdt_rate"`
}

func NewOrderDTO(m model.Order) *OrderDTO {
	var marketLink string
	if m.Exchange == model.ExchangeTypeBittrex {
		marketLink = "https://bittrex.com/Market/Index?MarketName=" + m.Market
	}
	return &OrderDTO{
		Market:      m.Market,
		MarketLink:  marketLink,
		Time:        m.Time.Unix(),
		BuyRate:     m.BuyRate,
		Amount:      m.Amount,
		SellNowRate: m.SellNowRate,
		USDTRate:    m.USDTRate,
	}
}
