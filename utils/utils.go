package utils

import (
	"bytes"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/shopspring/decimal"
)

type SpyLogger struct {
	sync.Mutex
	b bytes.Buffer
}

func (s *SpyLogger) String() string {
	s.Lock()
	defer s.Unlock()

	return s.b.String()
}

func (s *SpyLogger) Write(p []byte) (n int, err error) {
	s.Lock()
	defer s.Unlock()

	return s.b.Write(p)
}

// ExecuteConcurrently is helper function to execute set of tasks concurrently
// returns slice of errors, empty if no errors occurred
func ExecuteConcurrently(tasks []func() error) []error {
	var errors []error
	wg := sync.WaitGroup{}
	errChan := make(chan error, len(tasks))
	for _, task := range tasks {
		wg.Add(1)
		go func(task func() error) {
			defer wg.Done()
			err := task()
			if err != nil {
				errChan <- err
			}
		}(task)
	}
	wg.Wait()
	close(errChan)

	for err := range errChan {
		errors = append(errors, err)
	}

	return errors
}

// NewSpyLog creates fake logger containing produced output
func NewSpyLog() (*logrus.Entry, *SpyLogger) {
	var spyLogger SpyLogger

	logger := logrus.New()
	logger.Out = &spyLogger
	logger.Level = logrus.ErrorLevel

	return logrus.NewEntry(logger), &spyLogger
}

type devNullWriter struct {
}

func (devNullWriter) Write(p []byte) (n int, err error) {
	//do nothing
	return len(p), nil
}

// NewDevNullLog creates fake logger with nil output
func NewDevNullLog() *logrus.Entry {
	logger := logrus.New()
	logger.Out = devNullWriter{}
	logger.Level = logrus.ErrorLevel

	return logrus.NewEntry(logger)
}

// DecimalToFloatQuiet converts decimal to float64 in one line with exact ommitted
func DecimalToFloatQuiet(dec decimal.Decimal) float64 {
	f, _ := dec.Float64()
	return f
}
