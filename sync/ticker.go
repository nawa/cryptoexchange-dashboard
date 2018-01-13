package sync

import (
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"
)

const DefaultSyncTickerPeriod = time.Second

type Ticker struct {
	period      time.Duration
	timeTicker  *time.Ticker
	syncService Service
}

func NewSyncTicker(period time.Duration, service Service) *Ticker {
	return &Ticker{
		period:      period,
		syncService: service,
	}
}

func (s *Ticker) Start() error {
	if s.timeTicker != nil {
		return errors.New("ticker is already started")
	}
	ticker := time.NewTicker(s.period)
	go func() {
		for range ticker.C {
			err := s.syncService.Sync()
			if err != nil {
				log.Error(err)
			}
		}
	}()
	s.timeTicker = ticker
	return nil
}

func (s *Ticker) Stop() {
	if s.timeTicker == nil {
		return
	}

	s.timeTicker.Stop()
	s.timeTicker = nil
}
