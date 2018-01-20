package exchange

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/shopspring/decimal"
	"github.com/toorop/go-bittrex"

	"github.com/nawa/cryptoexchange-wallet-info/shared/model"
)

type BittrexExchange struct {
	bittrex         *bittrex.Bittrex
	marketSummaries []bittrex.MarketSummary
	syncTime        time.Time
}

func NewBittrexExchange(apiKey, apiSecret string) *BittrexExchange {
	bittrex := bittrex.New(apiKey, apiSecret)
	return &BittrexExchange{
		bittrex: bittrex,
	}
}

func (be *BittrexExchange) GetBalance() (*model.Balance, error) {
	var (
		balances []bittrex.Balance
		errCh    = make(chan error, 2)
	)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := be.updateMarketSummaries()
		if err != nil {
			errCh <- err
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		balances, err = be.bittrex.GetBalances()
		if err != nil {
			errCh <- err
		}
	}()

	wg.Wait()
	close(errCh)

	var err error
	for e := range errCh {
		err = multierror.Append(err, e)
	}

	if err != nil {
		return nil, err
	}

	result := &model.Balance{}
	for _, b := range balances {
		balance, _ := b.Balance.Float64()
		if balance > 0 {
			btcBalance, err := be.convertToBTC(b.Currency, balance)
			if err != nil {
				return nil, err
			}
			usdtBalance, err := be.convertToUSDT("BTC", btcBalance)
			if err != nil {
				return nil, err
			}
			btcRate, err := be.btcRate(b.Currency)
			if err != nil {
				return nil, err
			}

			result.Exchange = model.ExchangeTypeBittrex
			result.Time = be.syncTime
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

func (be *BittrexExchange) GetMarketInfo(market string) (*model.MarketInfo, error) {
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

func (be *BittrexExchange) updateMarketSummaries() (err error) {
	be.marketSummaries, err = be.bittrex.GetMarketSummaries()
	be.syncTime = time.Now().UTC()
	return
}

func (be *BittrexExchange) convertToBTC(currency string, amount float64) (float64, error) {
	return be.convertCurrency(currency, "BTC", amount)
}

func (be *BittrexExchange) convertToUSDT(currency string, amount float64) (float64, error) {
	return be.convertCurrency(currency, "USDT", amount)
}

func (be *BittrexExchange) btcRate(currency string) (float64, error) {
	return be.marketRate(currency, "BTC")
}

func (be *BittrexExchange) convertCurrency(fromCurrency, toCurrency string, amount float64) (float64, error) {
	if fromCurrency == toCurrency {
		return amount, nil
	}
	rate, err := be.marketRate(fromCurrency, toCurrency)
	if err != nil {
		return 0, err
	}
	return amount * rate, nil
}

func (be *BittrexExchange) marketRate(fromCurrency, toCurrency string) (float64, error) {
	if fromCurrency == toCurrency {
		return 1, nil
	}
	for _, market := range be.marketSummaries {
		if strings.ToUpper(market.MarketName) == strings.ToUpper(fmt.Sprintf("%s-%s", toCurrency, fromCurrency)) {
			last, _ := market.Last.Float64()
			return last, nil
		}
		if strings.ToUpper(market.MarketName) == strings.ToUpper(fmt.Sprintf("%s-%s", fromCurrency, toCurrency)) {
			reverse, _ := decimal.NewFromFloat(1).Div(market.Last).Float64()
			return reverse, nil
		}
	}
	return 0, fmt.Errorf("neither market '%s-%s' nor '%s-%s' found in markets", fromCurrency, toCurrency, toCurrency, fromCurrency)
}

func (be *BittrexExchange) Ping() error {
	_, err := be.bittrex.GetBalances()
	return err
}

func decToFloatQuiet(dec decimal.Decimal) float64 {
	f, _ := dec.Float64()
	return f
}
