package model

import "time"

type Order struct {
	//	|	Market	|	Time	|	BuyRate	| Can Sell=SellNowRate	|	Amount	|	Buy Price=BuyRate*Amount
	//	|	Sell Price=SellNowPrice*Amount	|	Profit = BuyPrice+BuyPrice*0.0025 - (Sell Price - Sell Price*0.0025)
	//	|	Profit BTC = Profit * BTCRate |	Profit USDT = Profit * USDTRate | Profit % = (Profit / Amount)*100
	Exchange    ExchangeType
	Market      string
	Time        time.Time
	BuyRate     float64
	Amount      float64
	SellNowRate float64
	USDTRate    float64
}
