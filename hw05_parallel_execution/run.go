package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}
	tasksCh := make(chan Task)

	go func() {
		for i, task := range tasks {
			if i == len(tasks)-1 {
				defer close(tasksCh)
			}

			tasksCh <- task
		}
	}()

	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			for {
				select {
				case v, ok := <-tasksCh:
					if !ok {
						wg.Done()

						return
					}

					v()
				}
			}
		}()
	}

	wg.Wait()

	return nil
}
