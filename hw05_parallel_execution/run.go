package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func watcher(tasks chan Task, errors chan bool, limit chan bool, killsignal chan bool, maxErrors int, wg *sync.WaitGroup, once *sync.Once) {
	for {
		select {
		case task := <-tasks:
			if len(errors) >= maxErrors {
				once.Do(func() { close(limit) })
			} else {
				err := task()
				if err != nil {
					errors <- true
				}
			}

			wg.Done()
		case <-killsignal:
			return
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	ch := make(chan Task, len(tasks))
	errors := make(chan bool, m+n)
	killsignal := make(chan bool)
	limit := make(chan bool)
	once := sync.Once{}
	wg := sync.WaitGroup{}

	defer close(errors)
	defer close(ch)
	defer close(killsignal)

	for i := 0; i < n; i++ {
		go watcher(ch, errors, limit, killsignal, m, &wg, &once)
	}

	for _, t := range tasks {
		wg.Add(1)
		ch <- t
	}

	wg.Wait()

	select {
	case <-limit:
		return ErrErrorsLimitExceeded
	default:
		return nil
	}
}
