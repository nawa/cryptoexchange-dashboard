package usecase

import (
	"encoding/json"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/nawa/cryptoexchange-wallet-info/storage"
)

type BalanceUsecase interface {
	StartSyncFromExchangePeriodically(period int) (*BalanceSynchronizer, error)
	SyncFromExchange() error
}

type balanceUsecase struct {
	exchange       storage.Exchange
	balanceStorage storage.BalanceStorage
}

func NewBalanceUsecase(exchange storage.Exchange, balanceStorage storage.BalanceStorage) BalanceUsecase {
	return &balanceUsecase{
		exchange:       exchange,
		balanceStorage: balanceStorage,
	}
}

func (b *balanceUsecase) StartSyncFromExchangePeriodically(period int) (*BalanceSynchronizer, error) {
	ticker := newSyncTicker(time.Second*time.Duration(period), b)
	err := ticker.start()
	return ticker, err
}

func (b *balanceUsecase) SyncFromExchange() error {
	balance, err := b.exchange.GetBalance()
	if err != nil {
		return err
	}

	err = b.balanceStorage.Save(storage.NewBalances(balance)...)
	if err != nil {
		return err
	}

	if log.GetLevel() >= log.DebugLevel {
		jsonBalance, err := json.MarshalIndent(balance, "", "  ")
		if err != nil {
			return err
		}
		log.WithField("balance", string(jsonBalance)).Debug()
	}
	return nil
}
