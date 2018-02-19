package ticker

import (
	"bytes"
	"testing"
	"time"

	"github.com/pkg/errors"
	assert "github.com/stretchr/testify/require"

	"github.com/nawa/cryptoexchange-dashboard/utils"
)

func TestTicker_Start(t *testing.T) {
	var i int
	ticker := NewTicker(time.Millisecond, func() error {
		i = i + 1
		return nil
	})
	ticker.log = utils.NewDevNullLog()
	err := ticker.Start()
	assert.NoError(t, err)
	time.Sleep(time.Millisecond * 10)

	assert.True(t, i > 2)
}

func TestTicker_Start_ErrorDoubleStart(t *testing.T) {
	ticker := NewTicker(time.Millisecond, func() error {
		return nil
	})
	ticker.log = utils.NewDevNullLog()

	err := ticker.Start()
	assert.NoError(t, err)

	err = ticker.Start()
	assert.Error(t, err)
}

func TestTicker_Start_ErrorFromTickF(t *testing.T) {
	ticker := NewTicker(time.Millisecond, func() error {
		return errors.New("some error")
	})
	var out *bytes.Buffer
	ticker.log, out = utils.NewSpyLog()

	err := ticker.Start()
	assert.NoError(t, err)

	time.Sleep(time.Millisecond * 10)

	assert.Contains(t, out.String(), "some error")
}

func TestTicker_Stop(t *testing.T) {
	var i int
	var ticker *Ticker
	ticker = NewTicker(time.Millisecond, func() error {
		ticker.Stop()
		//safe to stop twice
		ticker.Stop()
		ticker.Stop()
		i = i + 1
		return nil
	})
	ticker.log = utils.NewDevNullLog()
	err := ticker.Start()
	assert.NoError(t, err)
	time.Sleep(time.Millisecond * 10)

	assert.Equal(t, 1, i)
}
