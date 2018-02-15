package usecase

import (
	"encoding/json"
	"time"

	"github.com/nawa/cryptoexchange-dashboard/model"
	"github.com/nawa/cryptoexchange-dashboard/storage"

	"github.com/Sirupsen/logrus"
)

type BalanceUsecase interface {
	StartSyncFromExchangePeriodically(period int) (*BalanceSynchronizer, error)
	SyncFromExchange() error
	// All records from the last N hours
	FetchHourly(currency string, hours int) ([]model.CurrencyBalance, error)
	// Records from the last week with 5 min interval
	FetchWeekly(currency string) ([]model.CurrencyBalance, error)
	// Records from the last month with 1 hour interval
	FetchMonthly(currency string) ([]model.CurrencyBalance, error)
	//TODO // All records with 1 day???  interval
	FetchAll(currency string) ([]model.CurrencyBalance, error)
	// Get currency balances > 0
	GetActiveCurrencies() ([]model.CurrencyBalance, error)
}

type balanceUsecase struct {
	exchange       storage.Exchange
	balanceStorage storage.BalanceStorage
	log            *logrus.Entry
}

func NewBalanceUsecase(exchange storage.Exchange, balanceStorage storage.BalanceStorage) BalanceUsecase {
	log := logrus.WithField("component", "balanceUC")
	return &balanceUsecase{
		exchange:       exchange,
		balanceStorage: balanceStorage,
		log:            log,
	}
}

func (u *balanceUsecase) StartSyncFromExchangePeriodically(period int) (*BalanceSynchronizer, error) {
	ticker := newBalanceSynchronizer(time.Second*time.Duration(period), u)
	err := ticker.start()
	return ticker, err
}

func (u *balanceUsecase) SyncFromExchange() error {
	balance, err := u.exchange.GetBalance()
	if err != nil {
		u.log.WithField("method", "SyncFromExchange").WithError(err).Error()
		return err
	}

	err = u.balanceStorage.Save(storage.NewBalances(balance)...)
	if err != nil {
		u.log.WithField("method", "SyncFromExchange").WithError(err).Error()
		return err
	}

	if u.log.Level >= logrus.DebugLevel {
		jsonBalance, err := json.MarshalIndent(balance, "", "  ")
		if err != nil {
			u.log.WithField("method", "SyncFromExchange").WithError(err).Error()
			return err
		}
		u.log.WithField("balance", string(jsonBalance)).Debug("current balance")
	}
	return nil
}

func (u *balanceUsecase) FetchHourly(currency string, hours int) (balances []model.CurrencyBalance, err error) {
	stBalances, err := u.balanceStorage.FetchHourly(currency, hours)
	if err != nil {
		u.log.WithField("method", "FetchHourly").WithError(err).Error()
		return
	}
	for _, b := range stBalances {
		balances = append(balances, *b.ToModel())
	}
	return
}

func (u *balanceUsecase) FetchWeekly(currency string) (balances []model.CurrencyBalance, err error) {
	stBalances, err := u.balanceStorage.FetchWeekly(currency)
	if err != nil {
		u.log.WithField("method", "FetchWeekly").WithError(err).Error()
		return
	}
	for _, b := range stBalances {
		balances = append(balances, *b.ToModel())
	}
	return
}

func (u *balanceUsecase) FetchMonthly(currency string) (balances []model.CurrencyBalance, err error) {
	stBalances, err := u.balanceStorage.FetchMonthly(currency)
	if err != nil {
		u.log.WithField("method", "FetchMonthly").WithError(err).Error()
		return
	}
	for _, b := range stBalances {
		balances = append(balances, *b.ToModel())
	}
	return
}

func (u *balanceUsecase) FetchAll(currency string) (balances []model.CurrencyBalance, err error) {
	stBalances, err := u.balanceStorage.FetchAll(currency)
	if err != nil {
		u.log.WithField("method", "FetchAll").WithError(err).Error()
		return
	}
	for _, b := range stBalances {
		balances = append(balances, *b.ToModel())
	}
	return
}

func (u *balanceUsecase) GetActiveCurrencies() (balances []model.CurrencyBalance, err error) {
	stBalances, err := u.balanceStorage.GetActiveCurrencies()

	if err != nil {
		u.log.WithField("method", "FetchAll").WithError(err).Error()
		return
	}
	for _, b := range stBalances {
		balances = append(balances, *b.ToModel())
	}
	return
}
