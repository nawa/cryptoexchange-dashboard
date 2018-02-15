package utils

import (
	"sync"
)

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
