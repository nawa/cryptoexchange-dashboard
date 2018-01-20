package sync

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"

	"github.com/nawa/cryptoexchange-wallet-info/shared/exchange"
	"github.com/nawa/cryptoexchange-wallet-info/shared/storage"
)

type Service interface {
	Sync() error
}

type syncService struct {
	exchange       exchange.Exchange
	balanceStorage storage.BalanceStorage
}

func NewSyncService(exchange exchange.Exchange, balanceStorage storage.BalanceStorage) Service {
	return &syncService{
		exchange:       exchange,
		balanceStorage: balanceStorage,
	}
}

func (s *syncService) Sync() error {
	balance, err := s.exchange.GetBalance()
	if err != nil {
		return err
	}

	err = s.balanceStorage.Save(balance.ToStorageModel()...)
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
