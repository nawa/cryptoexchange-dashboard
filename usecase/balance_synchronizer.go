package usecase

import (
	"errors"
	"time"

	"github.com/Sirupsen/logrus"
)

const DefaultSyncTickerPeriod = time.Second

type BalanceSynchronizer struct {
	period         time.Duration
	timeTicker     *time.Ticker
	balanceUsecase BalanceUsecase
	log            *logrus.Entry
}

func newBalanceSynchronizer(period time.Duration, balanceUsecase BalanceUsecase) *BalanceSynchronizer {
	log := logrus.WithField("component", "balanceSynchronizer")
	return &BalanceSynchronizer{
		period:         period,
		balanceUsecase: balanceUsecase,
		log:            log,
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
				s.log.WithField("method", "start.ticker").WithError(err).Error()
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
