package exchange

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/nawa/cryptoexchange-dashboard/model"
	"github.com/nawa/cryptoexchange-dashboard/storage"
	"github.com/nawa/cryptoexchange-dashboard/utils"

	"github.com/Sirupsen/logrus"
	"github.com/hashicorp/go-multierror"
	"github.com/shopspring/decimal"
	"github.com/toorop/go-bittrex"
)

type bittrexExchange struct {
	bittrex *bittrex.Bittrex
	log     *logrus.Entry
}

type currencyConverter struct {
	marketSummaries []bittrex.MarketSummary
	syncTime        time.Time
}

func NewBittrexExchange(apiKey, apiSecret string) storage.Exchange {
	log := logrus.WithField("component", "BittrexExchange")
	bittrex := bittrex.New(apiKey, apiSecret)
	return &bittrexExchange{
		bittrex: bittrex,
		log:     log,
	}
}

func (be *bittrexExchange) GetBalance() (*model.Balance, error) {
	var (
		balances  []bittrex.Balance
		converter *currencyConverter
	)

	errs := utils.ExecuteConcurrently([]func() error{
		func() (err error) {
			converter, err = be.createCurrencyConverter()
			return
		},
		func() (err error) {
			balances, err = be.bittrex.GetBalances()
			return
		},
	})

	var err error
	for _, e := range errs {
		err = multierror.Append(err, e)
	}

	if err != nil {
		return nil, err
	}

	result := &model.Balance{}
	for _, b := range balances {
		balance, _ := b.Balance.Float64()
		if balance > 0 {
			btcBalance, err := converter.ConvertToBTC(b.Currency, balance)
			if err != nil {
				return nil, err
			}
			usdtBalance, err := converter.ConvertToUSDT("BTC", btcBalance)
			if err != nil {
				return nil, err
			}
			btcRate, err := converter.BtcRate(b.Currency)
			if err != nil {
				return nil, err
			}

			result.Exchange = model.ExchangeTypeBittrex
			result.Time = converter.syncTime
			result.Currencies = append(result.Currencies,
				model.CurrencyBalance{
					Currency:   b.Currency,
					Amount:     balance,
					BTCAmount:  btcBalance,
					USDTAmount: usdtBalance,
					BTCRate:    btcRate,
				})
			result.BTCAmount = result.BTCAmount + btcBalance
			result.USDTAmount = result.USDTAmount + usdtBalance
		}
	}
	return result, nil
}

func (be *bittrexExchange) GetMarketInfo(market string) (*model.MarketInfo, error) {
	marketSummary, err := be.bittrex.GetMarketSummary(market)
	if err != nil {
		return nil, err
	}

	if len(marketSummary) == 0 {
		return nil, errors.New("got empty marketSummary list")
	}

	return &model.MarketInfo{
		MarketName: marketSummary[0].MarketName,
		Last:       decToFloatQuiet(marketSummary[0].Last),
		Bid:        decToFloatQuiet(marketSummary[0].Bid),
		Ask:        decToFloatQuiet(marketSummary[0].Ask),
		High:       decToFloatQuiet(marketSummary[0].High),
		Low:        decToFloatQuiet(marketSummary[0].Low),
	}, nil
}

func (be *bittrexExchange) GetOrders() ([]model.Order, error) {
	var (
		orders    []bittrex.Order
		converter *currencyConverter
	)

	errs := utils.ExecuteConcurrently([]func() error{
		func() (err error) {
			converter, err = be.createCurrencyConverter()
			return
		},
		func() (err error) {
			orders, err = be.bittrex.GetOrderHistory("all")
			return
		},
	})

	var err error
	for _, e := range errs {
		err = multierror.Append(err, e)
	}

	if err != nil {
		return nil, err
	}

	sold := make(map[string]bool)
	filtered := make([]bittrex.Order, 0)

	// leave only last buy orders in each active market
	for _, order := range orders {
		if !sold[order.Exchange] {
			if order.OrderType == "LIMIT_BUY" {
				filtered = append(filtered, order)
			} else {
				sold[order.Exchange] = true
			}
		}
	}

	return be.convertOrders(filtered, converter), nil
}

func (be *bittrexExchange) convertOrders(bittrexOrders []bittrex.Order, converter *currencyConverter) (orders []model.Order) {
	for _, order := range bittrexOrders {
		toFrom := strings.Split(order.Exchange, "-")
		if len(toFrom) != 2 {
			be.log.WithField("method", "convertOrders").Warnf("exchange name can'be parsed to from-to format - %s", order.Exchange)
			continue
		}

		_, bidRate, _, err := converter.MarketRate(toFrom[1], toFrom[0])
		if err != nil {
			be.log.WithField("method", "convertOrders").Warnf("market rate can't be found")
			continue
		}

		usdtRate, err := converter.ConvertToUSDT(toFrom[0], 1)
		if err != nil {
			be.log.WithField("method", "convertOrders").Warnf("market convert to USDT")
			continue
		}
		quanity, _ := order.Quantity.Float64()
		buyRate, _ := order.Price.Div(order.Quantity).Float64()

		orders = append(orders, model.Order{
			Exchange:    model.ExchangeTypeBittrex,
			Market:      order.Exchange,
			Time:        order.TimeStamp.Time,
			Amount:      quanity,
			BuyRate:     buyRate,
			SellNowRate: bidRate,
			USDTRate:    usdtRate,
		})
	}
	return
}

func (be *bittrexExchange) Ping() error {
	_, err := be.bittrex.GetBalances()
	return err
}

func (be *bittrexExchange) createCurrencyConverter() (*currencyConverter, error) {
	marketSummaries, err := be.bittrex.GetMarketSummaries()
	if err != nil {
		return nil, err
	}

	return &currencyConverter{
		marketSummaries: marketSummaries,
		syncTime:        time.Now().UTC(),
	}, nil
}

func (c *currencyConverter) ConvertToBTC(currency string, amount float64) (float64, error) {
	return c.ConvertCurrency(currency, "BTC", amount)
}

func (c *currencyConverter) ConvertToUSDT(currency string, amount float64) (float64, error) {
	return c.ConvertCurrency(currency, "USDT", amount)
}

func (c *currencyConverter) BtcRate(currency string) (float64, error) {
	last, _, _, err := c.MarketRate(currency, "BTC")
	return last, err
}

func (c *currencyConverter) ConvertCurrency(fromCurrency, toCurrency string, amount float64) (float64, error) {
	if fromCurrency == toCurrency {
		return amount, nil
	}
	last, _, _, err := c.MarketRate(fromCurrency, toCurrency)
	if err != nil {
		return 0, err
	}
	return amount * last, nil
}

func (c *currencyConverter) MarketRate(fromCurrency, toCurrency string) (last float64, bid float64, ask float64, err error) {
	if fromCurrency == toCurrency {
		return 1, 1, 1, nil
	}
	for _, market := range c.marketSummaries {
		if strings.ToUpper(market.MarketName) == strings.ToUpper(fmt.Sprintf("%s-%s", toCurrency, fromCurrency)) {
			last, _ = market.Last.Float64()
			bid, _ = market.Bid.Float64()
			ask, _ = market.Ask.Float64()
			return
		}
		if strings.ToUpper(market.MarketName) == strings.ToUpper(fmt.Sprintf("%s-%s", fromCurrency, toCurrency)) {
			last, _ = decimal.NewFromFloat(1).Div(market.Last).Float64()
			bid, _ = decimal.NewFromFloat(1).Div(market.Bid).Float64()
			ask, _ = decimal.NewFromFloat(1).Div(market.Ask).Float64()
			return
		}
	}
	return 0, 0, 0, fmt.Errorf("neither market '%s-%s' nor '%s-%s' found in markets", fromCurrency, toCurrency, toCurrency, fromCurrency)
}

func decToFloatQuiet(dec decimal.Decimal) float64 {
	f, _ := dec.Float64()
	return f
}
