package usecase

import (
	"encoding/json"
	"time"

	"github.com/nawa/cryptoexchange-dashboard/domain"
	"github.com/nawa/cryptoexchange-dashboard/usecase/ticker"

	"github.com/nawa/cryptoexchange-dashboard/storage"

	"github.com/Sirupsen/logrus"
)

type BalanceUsecases interface {
	StartSyncFromExchangePeriodically(period time.Duration) (stop func(), err error)
	SyncFromExchange() error
	// All records from the last N hours
	FetchHourly(currency string, hours int) ([]domain.Balance, error)
	// Records from the last week with 5 min interval
	FetchWeekly(currency string) ([]domain.Balance, error)
	// Records from the last month with 1 hour interval
	FetchMonthly(currency string) ([]domain.Balance, error)
	//TODO // All records with 1 day???  interval
	FetchAll(currency string) ([]domain.Balance, error)
	// Get currency balances > 0
	GetActiveCurrencies() ([]domain.Balance, error)
}

type balanceUsecases struct {
	exchange       storage.Exchange
	balanceStorage storage.BalanceStorage
	log            *logrus.Entry
}

func NewBalanceUsecase(exchange storage.Exchange, balanceStorage storage.BalanceStorage) BalanceUsecases {
	log := logrus.WithField("component", "balanceUC")
	return &balanceUsecases{
		exchange:       exchange,
		balanceStorage: balanceStorage,
		log:            log,
	}
}

func (u *balanceUsecases) StartSyncFromExchangePeriodically(period time.Duration) (stop func(), err error) {
	ticker := ticker.NewTicker(period, u.SyncFromExchange)
	err = ticker.Start()
	if err != nil {
		return nil, err
	}

	return func() {
		ticker.Stop()
	}, err
}

func (u *balanceUsecases) SyncFromExchange() error {
	balances, err := u.exchange.GetBalance()
	if err != nil {
		u.log.WithField("method", "SyncFromExchange").WithError(err).Error()
		return err
	}

	if len(balances) == 0 {
		return nil
	}

	total := domain.Balance{
		Currency: "total",
		//TODO fix me for multiple exchanges
		Exchange: balances[0].Exchange,
		Time:     balances[0].Time,
	}

	for _, b := range balances {
		total.BTCAmount += b.BTCAmount
		total.USDTAmount += b.USDTAmount
	}

	balances = append(balances, total)

	err = u.balanceStorage.Save(balances...)
	if err != nil {
		u.log.WithField("method", "SyncFromExchange").WithError(err).Error()
		return err
	}

	if u.log.Level >= logrus.DebugLevel {
		jsonBalances, err := json.MarshalIndent(balances, "", "  ")
		if err != nil {
			u.log.WithField("method", "SyncFromExchange").WithError(err).Error()
			return err
		}
		u.log.WithField("balance", string(jsonBalances)).Debug("current balance")
	}
	return nil
}

func (u *balanceUsecases) FetchHourly(currency string, hours int) ([]domain.Balance, error) {
	balances, err := u.balanceStorage.FetchHourly(currency, hours)
	if err != nil {
		u.log.WithField("method", "FetchHourly").WithError(err).Error()
		return nil, err
	}

	return balances, nil
}

func (u *balanceUsecases) FetchWeekly(currency string) ([]domain.Balance, error) {
	balances, err := u.balanceStorage.FetchWeekly(currency)
	if err != nil {
		u.log.WithField("method", "FetchWeekly").WithError(err).Error()
		return nil, err
	}

	return balances, nil
}

func (u *balanceUsecases) FetchMonthly(currency string) ([]domain.Balance, error) {
	balances, err := u.balanceStorage.FetchMonthly(currency)
	if err != nil {
		u.log.WithField("method", "FetchMonthly").WithError(err).Error()
		return nil, err
	}

	return balances, nil
}

func (u *balanceUsecases) FetchAll(currency string) ([]domain.Balance, error) {
	balances, err := u.balanceStorage.FetchAll(currency)
	if err != nil {
		u.log.WithField("method", "FetchAll").WithError(err).Error()
		return nil, err
	}

	return balances, nil
}

func (u *balanceUsecases) GetActiveCurrencies() ([]domain.Balance, error) {
	balances, err := u.balanceStorage.GetActiveCurrencies()
	if err != nil {
		u.log.WithField("method", "GetActiveCurrencies").WithError(err).Error()
		return nil, err
	}

	return balances, nil
}
