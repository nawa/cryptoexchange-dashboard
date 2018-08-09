package ticker

import (
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

type TickF func() error

type Ticker struct {
	period     time.Duration
	timeTicker *time.Ticker
	tickerF    TickF
	log        *logrus.Entry
	lock       sync.Mutex
}

func NewTicker(period time.Duration, tickerF TickF) *Ticker {
	log := logrus.WithField("component", "Ticker")
	return &Ticker{
		period:  period,
		tickerF: tickerF,
		log:     log,
	}
}

func (t *Ticker) Start() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.timeTicker != nil {
		return errors.New("ticker is already started")
	}
	ticker := time.NewTicker(t.period)
	go func() {
		for range ticker.C {
			err := t.tickerF()
			if err != nil {
				t.log.WithField("method", "start.ticker").WithError(err).Error()
			}
		}
	}()
	t.timeTicker = ticker
	return nil
}

func (t *Ticker) Stop() {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.timeTicker == nil {
		return
	}

	t.timeTicker.Stop()
	t.timeTicker = nil
}
