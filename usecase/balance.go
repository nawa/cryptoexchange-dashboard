package usecase

import (
	"encoding/json"
	"time"

	"github.com/nawa/cryptoexchange-wallet-info/model"
	"github.com/nawa/cryptoexchange-wallet-info/storage"

	"github.com/Sirupsen/logrus"
)

type BalanceUsecase interface {
	StartSyncFromExchangePeriodically(period int) (*BalanceSynchronizer, error)
	SyncFromExchange() error
	FetchDaily(currency string) ([]model.CurrencyBalance, error)
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

func (u *balanceUsecase) FetchDaily(currency string) (balances []model.CurrencyBalance, err error) {
	stBalances, err := u.balanceStorage.FetchDaily(currency)
	if err != nil {
		u.log.WithField("method", "FetchDaily").WithError(err).Error()
		return
	}
	for _, b := range stBalances {
		balances = append(balances, *b.ToModel())
	}
	return
}
