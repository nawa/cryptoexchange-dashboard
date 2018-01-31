package usecase

import (
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"
)

const DefaultSyncTickerPeriod = time.Second

type BalanceSynchronizer struct {
	period         time.Duration
	timeTicker     *time.Ticker
	balanceUsecase BalanceUsecase
}

func newSyncTicker(period time.Duration, balanceUsecase BalanceUsecase) *BalanceSynchronizer {
	return &BalanceSynchronizer{
		period:         period,
		balanceUsecase: balanceUsecase,
	}
}

func (s *BalanceSynchronizer) start() error {
	if s.timeTicker != nil {
		return errors.New("ticker is already started")
	}
	ticker := time.NewTicker(s.period)
	go func() {
		for range ticker.C {
			err := s.balanceUsecase.SyncFromExchange()
			if err != nil {
				log.Error(err)
			}
		}
	}()
	s.timeTicker = ticker
	return nil
}

func (s *BalanceSynchronizer) Stop() {
	if s.timeTicker == nil {
		return
	}

	s.timeTicker.Stop()
	s.timeTicker = nil
}
