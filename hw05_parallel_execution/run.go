package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrNoWorkers           = errors.New("no workers found")
	ErrMaxErrorsIsNotValid = errors.New("max errors number is too small")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}
	tasksCh := make(chan Task)
	var errors int64
	var err error

	if n <= 0 {
		err = ErrNoWorkers

		return err
	}

	if m <= 0 {
		err = ErrMaxErrorsIsNotValid

		return err
	}

	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			for v := range tasksCh {
				if v == nil {
					continue
				}

				err := v()
				if err != nil {
					atomic.AddInt64(&errors, 1)
				}
			}

			wg.Done()
		}()
	}

	for _, task := range tasks {
		if atomic.LoadInt64(&errors) >= int64(m) {
			err = ErrErrorsLimitExceeded

			break
		}

		tasksCh <- task
	}

	close(tasksCh)

	wg.Wait()

	return err
}
