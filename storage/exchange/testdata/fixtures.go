package testdata

import (
	"encoding/json"
	"time"

	"github.com/nawa/cryptoexchange-dashboard/utils"

	"github.com/nawa/cryptoexchange-dashboard/model"
	"github.com/shopspring/decimal"
	bittrex "github.com/toorop/go-bittrex"
)

type BittrexJSONResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

var (
	usdtMarketSummary = bittrex.MarketSummary{
		MarketName: "BTC-USDT",
		Last:       decimal.NewFromFloat(130),
		Bid:        decimal.NewFromFloat(140),
		Ask:        decimal.NewFromFloat(150),
	}
)

func BittrexMarketSummaries() []bittrex.MarketSummary {
	return []bittrex.MarketSummary{
		{
			MarketName: "BTC-CUR1",
			Last:       decimal.NewFromFloat(10),
			Bid:        decimal.NewFromFloat(20),
			Ask:        decimal.NewFromFloat(30),
		},
		{
			MarketName: "BTC-CUR2",
			Last:       decimal.NewFromFloat(40),
			Bid:        decimal.NewFromFloat(50),
			Ask:        decimal.NewFromFloat(60),
		},
		{
			MarketName: "CUR1-CUR2",
			Last:       decimal.NewFromFloat(100),
			Bid:        decimal.NewFromFloat(110),
			Ask:        decimal.NewFromFloat(120),
		},
		usdtMarketSummary,
	}
}

func BittrexBalances() []bittrex.Balance {
	return []bittrex.Balance{
		{
			Currency: "BTC",
			Balance:  decimal.NewFromFloat(1000),
		},
		{
			Currency: "CUR1",
			Balance:  decimal.NewFromFloat(2000),
		},
		{
			Currency: "CUR2",
			Balance:  decimal.NewFromFloat(3000),
		},
	}
}

func ModelBalance() *model.Balance {
	currencyBalances := []model.CurrencyBalance{
		{
			Amount:     1000,
			BTCAmount:  1000,
			Currency:   "BTC",
			USDTAmount: utils.DecimalToFloatQuiet(decimal.NewFromFloat(1000).Mul(decimal.NewFromFloat(1).Div(usdtMarketSummary.Last))),
		},
		{
			Amount:     2000,
			BTCAmount:  20000,
			Currency:   "CUR1",
			USDTAmount: utils.DecimalToFloatQuiet(decimal.NewFromFloat(20000).Mul(decimal.NewFromFloat(1).Div(usdtMarketSummary.Last))),
		},
		{
			Amount:     3000,
			BTCAmount:  120000,
			Currency:   "CUR2",
			USDTAmount: utils.DecimalToFloatQuiet(decimal.NewFromFloat(120000).Mul(decimal.NewFromFloat(1).Div(usdtMarketSummary.Last))),
		},
	}

	return &model.Balance{
		Currencies: currencyBalances,
		Exchange:   model.ExchangeTypeBittrex,
		BTCAmount:  currencyBalances[0].BTCAmount + currencyBalances[1].BTCAmount + currencyBalances[2].BTCAmount,
		USDTAmount: currencyBalances[0].USDTAmount + currencyBalances[1].USDTAmount + currencyBalances[2].USDTAmount,
	}
}

func ModelMarketInfo() *model.MarketInfo {
	marketSummary := BittrexMarketSummaries()[0]
	return &model.MarketInfo{
		MarketName: marketSummary.MarketName,
		Last:       utils.DecimalToFloatQuiet(marketSummary.Last),
		Bid:        utils.DecimalToFloatQuiet(marketSummary.Bid),
		Ask:        utils.DecimalToFloatQuiet(marketSummary.Ask),
		High:       utils.DecimalToFloatQuiet(marketSummary.High),
		Low:        utils.DecimalToFloatQuiet(marketSummary.Low),
	}
}
func BittrexOrders() []bittrex.Order {
	orders := []bittrex.Order{
		{
			Exchange:  "BTC-CUR1",
			OrderType: "LIMIT_BUY",
			Quantity:  decimal.NewFromFloat(100),
			Price:     decimal.NewFromFloat(10),
		},
		{
			Exchange:  "BTC-CUR2",
			OrderType: "LIMIT_BUY",
			Quantity:  decimal.NewFromFloat(200),
			Price:     decimal.NewFromFloat(20),
		},
		{
			Exchange:  "BTC-CUR1",
			OrderType: "LIMIT_BUY",
			Quantity:  decimal.NewFromFloat(300),
			Price:     decimal.NewFromFloat(30),
		},
		{
			Exchange:  "BTC-CUR1",
			OrderType: "LIMIT_SELL",
			Quantity:  decimal.NewFromFloat(800),
			Price:     decimal.NewFromFloat(100),
		},
		{
			Exchange:  "BTC-CUR1",
			OrderType: "LIMIT_BUY",
			Quantity:  decimal.NewFromFloat(1000),
			Price:     decimal.NewFromFloat(100),
		},
		{
			Exchange:  "BTC-CUR2",
			OrderType: "LIMIT_BUY",
			Quantity:  decimal.NewFromFloat(300),
			Price:     decimal.NewFromFloat(50),
		},
		{
			Exchange:  "BTC-CUR2",
			OrderType: "LIMIT_SELL",
			Quantity:  decimal.NewFromFloat(100),
			Price:     decimal.NewFromFloat(20),
		},
		{
			Exchange:  "BTC-CUR2",
			OrderType: "LIMIT_BUY",
			Quantity:  decimal.NewFromFloat(300),
			Price:     decimal.NewFromFloat(50),
		},
		{ //Bad order
			Exchange:  "CUR5-CUR6",
			OrderType: "LIMIT_BUY",
			Quantity:  decimal.NewFromFloat(300),
			Price:     decimal.NewFromFloat(50),
		},
		{ //Bad order
			Exchange:  "ABCDE",
			OrderType: "LIMIT_BUY",
			Quantity:  decimal.NewFromFloat(300),
			Price:     decimal.NewFromFloat(50),
		},
		{ //Bad order - can't convert to USDT
			Exchange:  "CUR1-CUR2",
			OrderType: "LIMIT_BUY",
			Quantity:  decimal.NewFromFloat(300),
			Price:     decimal.NewFromFloat(50),
		},
	}
	basetime := time.Unix(0, 0).UTC().Add(time.Hour * 24 * 1000)

	orders[0].TimeStamp.Time = basetime
	orders[1].TimeStamp.Time = basetime.Add(-time.Hour)
	orders[2].TimeStamp.Time = basetime.Add(-time.Hour * 2)
	orders[3].TimeStamp.Time = basetime.Add(-time.Hour * 3)
	orders[4].TimeStamp.Time = basetime.Add(-time.Hour * 4)
	orders[5].TimeStamp.Time = basetime.Add(-time.Hour * 5)
	return orders
}

func ModelOrders() []model.Order {
	basetime := time.Unix(0, 0).UTC().Add(time.Hour * 24 * 1000)
	s := basetime.Format(bittrex.TIME_FORMAT)
	basetimeAfterBittrexSerialization, err := time.Parse(bittrex.TIME_FORMAT, s)
	if err != nil {
		panic(err)
	}
	bittrexOrders := BittrexOrders()
	bittrexMarketSummaries := BittrexMarketSummaries()
	return []model.Order{
		{
			Market:      "BTC-CUR1",
			Exchange:    model.ExchangeTypeBittrex,
			Time:        basetimeAfterBittrexSerialization,
			Amount:      utils.DecimalToFloatQuiet(bittrexOrders[0].Quantity),
			BuyRate:     utils.DecimalToFloatQuiet(bittrexOrders[0].Price.Div(bittrexOrders[0].Quantity)),
			SellNowRate: utils.DecimalToFloatQuiet(bittrexMarketSummaries[0].Bid),
			USDTRate:    utils.DecimalToFloatQuiet(decimal.NewFromFloat(1).Div(usdtMarketSummary.Last)),
		},
		{
			Market:      "BTC-CUR2",
			Exchange:    model.ExchangeTypeBittrex,
			Time:        basetimeAfterBittrexSerialization.Add(-time.Hour),
			Amount:      utils.DecimalToFloatQuiet(bittrexOrders[1].Quantity),
			BuyRate:     utils.DecimalToFloatQuiet(bittrexOrders[1].Price.Div(bittrexOrders[1].Quantity)),
			SellNowRate: utils.DecimalToFloatQuiet(bittrexMarketSummaries[1].Bid),
			USDTRate:    utils.DecimalToFloatQuiet(decimal.NewFromFloat(1).Div(usdtMarketSummary.Last)),
		},
		{
			Market:      "BTC-CUR1",
			Exchange:    model.ExchangeTypeBittrex,
			Time:        basetimeAfterBittrexSerialization.Add(-time.Hour * 2),
			Amount:      utils.DecimalToFloatQuiet(bittrexOrders[2].Quantity),
			BuyRate:     utils.DecimalToFloatQuiet(bittrexOrders[2].Price.Div(bittrexOrders[2].Quantity)),
			SellNowRate: utils.DecimalToFloatQuiet(bittrexMarketSummaries[0].Bid),
			USDTRate:    utils.DecimalToFloatQuiet(decimal.NewFromFloat(1).Div(usdtMarketSummary.Last)),
		},
		{
			Market:      "BTC-CUR2",
			Exchange:    model.ExchangeTypeBittrex,
			Time:        basetimeAfterBittrexSerialization.Add(-time.Hour * 5),
			Amount:      utils.DecimalToFloatQuiet(bittrexOrders[5].Quantity),
			BuyRate:     utils.DecimalToFloatQuiet(bittrexOrders[5].Price.Div(bittrexOrders[5].Quantity)),
			SellNowRate: utils.DecimalToFloatQuiet(bittrexMarketSummaries[1].Bid),
			USDTRate:    utils.DecimalToFloatQuiet(decimal.NewFromFloat(1).Div(usdtMarketSummary.Last)),
		},
	}
}

func BittrexResponseSuccess(result interface{}) interface{} {
	marshalled, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	return &BittrexJSONResponse{
		Success: true,
		Message: "OK",
		Result:  marshalled,
	}
}

func BittrexResponseFailure() *BittrexJSONResponse {
	return &BittrexJSONResponse{
		Success: false,
		Message: "NOT OK",
		Result:  []byte{},
	}
}
