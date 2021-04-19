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

	for i, task := range tasks {
		go func(t Task, index int) {
			if index == len(tasks)-1 {
				defer close(tasksCh)
			}

			tasksCh <- t
		}(task, i)
	}

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
