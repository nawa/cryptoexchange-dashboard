package exchange

import (
	"fmt"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/toorop/go-bittrex"

	"github.com/nawa/cryptoexchange-dashboard/domain"
	"github.com/nawa/cryptoexchange-dashboard/storage"
	"github.com/nawa/cryptoexchange-dashboard/utils"
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

func (be *bittrexExchange) GetBalance() ([]domain.Balance, error) {
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

	var result []domain.Balance
	for _, b := range balances {
		if b.Balance.GreaterThan(decimal.NewFromFloat(0)) {
			btcBalance, err := converter.ConvertToBTC(b.Currency, b.Balance)
			if err != nil {
				return nil, err
			}
			usdtBalance, err := converter.ConvertToUSDT("BTC", btcBalance)
			if err != nil {
				return nil, err
			}

			result = append(result, domain.Balance{
				Exchange:   domain.ExchangeTypeBittrex,
				Currency:   b.Currency,
				Amount:     utils.DecimalToFloatQuiet(b.Balance),
				BTCAmount:  utils.DecimalToFloatQuiet(btcBalance),
				USDTAmount: utils.DecimalToFloatQuiet(usdtBalance),
				Time:       converter.syncTime,
			})
		}
	}
	return result, nil
}

func (be *bittrexExchange) GetMarketInfo(market string) (*domain.MarketInfo, error) {
	marketSummary, err := be.bittrex.GetMarketSummary(market)
	if err != nil {
		return nil, err
	}

	if len(marketSummary) == 0 {
		return nil, errors.New("got empty marketSummary list")
	}

	return &domain.MarketInfo{
		MarketName: marketSummary[0].MarketName,
		Last:       utils.DecimalToFloatQuiet(marketSummary[0].Last),
		Bid:        utils.DecimalToFloatQuiet(marketSummary[0].Bid),
		Ask:        utils.DecimalToFloatQuiet(marketSummary[0].Ask),
		High:       utils.DecimalToFloatQuiet(marketSummary[0].High),
		Low:        utils.DecimalToFloatQuiet(marketSummary[0].Low),
	}, nil
}

func (be *bittrexExchange) GetOrders() ([]domain.Order, error) {
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

func (be *bittrexExchange) convertOrders(bittrexOrders []bittrex.Order, converter *currencyConverter) []domain.Order {
	orders := []domain.Order{} //don't change me
	for _, order := range bittrexOrders {
		toFrom := strings.Split(order.Exchange, "-")
		if len(toFrom) != 2 {
			be.log.WithField("method", "convertOrders").Warnf("exchange name can't be parsed to from-to format - %s", order.Exchange)
			continue
		}

		_, bidRate, _, err := converter.MarketRate(toFrom[1], toFrom[0])
		if err != nil {
			be.log.WithField("method", "convertOrders").Warnf("market rate can't be found")
			continue
		}

		usdtRate, err := converter.ConvertToUSDT(toFrom[0], decimal.NewFromFloat(1))
		if err != nil {
			be.log.WithField("method", "convertOrders").Warnf("market convert to USDT")
			continue
		}

		orders = append(orders, domain.Order{
			Exchange:    domain.ExchangeTypeBittrex,
			Market:      order.Exchange,
			Time:        order.TimeStamp.Time,
			Amount:      utils.DecimalToFloatQuiet(order.Quantity),
			BuyRate:     utils.DecimalToFloatQuiet(order.Price.Div(order.Quantity)),
			SellNowRate: utils.DecimalToFloatQuiet(bidRate),
			USDTRate:    utils.DecimalToFloatQuiet(usdtRate),
		})
	}
	return orders
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

func (c *currencyConverter) ConvertToBTC(currency string, amount decimal.Decimal) (decimal.Decimal, error) {
	return c.ConvertCurrency(currency, "BTC", amount)
}

func (c *currencyConverter) ConvertToUSDT(currency string, amount decimal.Decimal) (decimal.Decimal, error) {
	return c.ConvertCurrency(currency, "USDT", amount)
}

func (c *currencyConverter) ConvertCurrency(fromCurrency, toCurrency string, amount decimal.Decimal) (decimal.Decimal, error) {
	last, _, _, err := c.MarketRate(fromCurrency, toCurrency)
	if err != nil {
		return decimal.Decimal{}, err
	}
	return amount.Mul(last), nil
}

func (c *currencyConverter) MarketRate(fromCurrency, toCurrency string) (last decimal.Decimal, bid decimal.Decimal, ask decimal.Decimal, err error) {
	if fromCurrency == toCurrency {
		return decimal.NewFromFloat(1), decimal.NewFromFloat(1), decimal.NewFromFloat(1), nil
	}
	for _, market := range c.marketSummaries {
		if strings.ToUpper(market.MarketName) == strings.ToUpper(fmt.Sprintf("%s-%s", toCurrency, fromCurrency)) {
			last = market.Last
			bid = market.Bid
			ask = market.Ask
			return
		}
		if strings.ToUpper(market.MarketName) == strings.ToUpper(fmt.Sprintf("%s-%s", fromCurrency, toCurrency)) {
			last = decimal.NewFromFloat(1).Div(market.Last)
			bid = decimal.NewFromFloat(1).Div(market.Bid)
			ask = decimal.NewFromFloat(1).Div(market.Ask)
			return
		}
	}
	err = fmt.Errorf("neither market '%s-%s' nor '%s-%s' found in markets", fromCurrency, toCurrency, toCurrency, fromCurrency)
	return
}

// func (c *currencyConverter) MarketRateWorkaround(fromCurrency /* CUR3 */, toCurrency /* BTC */ string) (last decimal.Decimal, bid decimal.Decimal, ask decimal.Decimal, err error) {
// 	for _, market := range c.marketSummaries {
// 		toFrom := strings.Split(strings.ToUpper(market.MarketName), "-")
// 		if len(toFrom) != 2 {
// 			continue
// 		}
// 		if toFrom[0] == fromCurrency { /* CUR3-ETH */
// 			last, bid, ask, err := c.MarketRate(toFrom[1] /* ETH */, toCurrency /* BTC */)
// 			if err == nil {
// 				return last.Mul(market.Last), bid.Mul(market.Bid), ask.Mul(market.Ask), nil
// 			}
// 		}

// 		if toFrom[1] == fromCurrency { /* ETH-CUR3 */
// 			last, bid, ask, err := c.MarketRate(toFrom[0] /* ETH */, toCurrency /* BTC */)
// 			if err == nil {
// 				return last.Mul(market.Last), bid.Mul(market.Bid), ask.Mul(market.Ask), nil
// 			}
// 		}
// 	}
// 	err = fmt.Errorf("neither market '%s-%s' nor '%s-%s' found in markets", fromCurrency, toCurrency, toCurrency, fromCurrency)
// 	return
// }
