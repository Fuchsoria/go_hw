package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func creator(tasks []Task, killSignal <-chan struct{}) <-chan Task {
	results := make(chan Task)

	go func() {
		defer close(results)

		for _, task := range tasks {
			select {
			case <-killSignal:
				break
			case results <- task:
			}
		}
	}()

	return results
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var errorsCount int64 = 0
	killSignal := make(chan struct{})
	done := make(chan struct{})
	limit := make(chan struct{})
	once := sync.Once{}
	onceLimit := sync.Once{}
	wg := sync.WaitGroup{}

	defer close(killSignal)

	worker := func(results <-chan Task) {
		defer wg.Done()

		for {
			select {
			case result, ok := <-results:
				if !ok {
					once.Do(func() {
						close(done)
					})
					return
				}
				if err := result(); err != nil {
					atomic.AddInt64(&errorsCount, 1)
				}
				if atomic.LoadInt64(&errorsCount) >= int64(m) {
					onceLimit.Do((func() {
						close(limit)
					}))
					return
				}
			case <-killSignal:
				return
			}
		}
	}

	tasksCh := creator(tasks, killSignal)

	for i := 0; i < n; i++ {
		wg.Add(1)

		go worker(tasksCh)
	}

	wg.Wait()

	for {
		select {
		case <-done:
			return nil
		case <-limit:
			return ErrErrorsLimitExceeded
		}
	}
}
