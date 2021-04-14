package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func watcher(tasks chan Task, errorsCount *int64, limit chan struct{}, killsignal chan struct{}, maxErrors int, wg *sync.WaitGroup, once *sync.Once, once2 *sync.Once) {
	closeSignal := func() { once.Do(func() { close(killsignal) }) }
	closeLimit := func() { once2.Do(func() { close(limit) }) }

	for {
		select {
		case task := <-tasks:
			wg.Add(1)
			err := task()
			if err != nil {
				atomic.AddInt64(errorsCount, 1)
			}
			wg.Done()

			if errors := int(*errorsCount); errors >= maxErrors {
				closeSignal()
				closeLimit()

				return
			}

			if len(tasks) == 0 {
				closeSignal()

				return
			}

		case <-killsignal:
			return
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	ch := make(chan Task, len(tasks))
	killsignal := make(chan struct{})
	limit := make(chan struct{})

	var errorsCount int64

	once := sync.Once{}
	once2 := sync.Once{}
	wg := sync.WaitGroup{}

	defer close(ch)

	for i := 0; i < n; i++ {
		go watcher(ch, &errorsCount, limit, killsignal, m, &wg, &once, &once2)
	}

	for _, t := range tasks {
		ch <- t
	}

	<-killsignal

	wg.Wait()

	select {
	case <-limit:
		return ErrErrorsLimitExceeded
	default:
		return nil
	}
}
