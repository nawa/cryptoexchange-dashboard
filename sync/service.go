package sync

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/nawa/cryptoexchange-wallet-info/sync/exchange"
)

type Service interface {
	Sync() error
}

type syncService struct {
	exchange exchange.Exchange
}

func NewSyncService(exchange exchange.Exchange) Service {
	return &syncService{
		exchange: exchange,
	}
}

func (s *syncService) Sync() error {
	balance, err := s.exchange.GetBalance()
	if err != nil {
		return err
	}
	jsonBalance, err := json.MarshalIndent(balance, "", "  ")
	if err != nil {
		return err
	}
	log.WithField("balance", string(jsonBalance)).Debug()
	return nil
}
