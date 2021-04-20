package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("Run with nil values", func(t *testing.T) {
		tasks := []Task{}
		maxTasks := 50
		workersCount := 10
		maxErrorsCount := 1
		var nonNilTasks int64
		var doneTasks int64

		randBool := func() bool {
			return rand.Intn(2) == 1
		}

		for i := 0; i < maxTasks; i++ {
			if nonNil := randBool(); nonNil {
				nonNilTasks++

				tasks = append(tasks, func() error {
					atomic.AddInt64(&doneTasks, 1)

					return nil
				})
			} else {
				tasks = append(tasks, nil)
			}
		}

		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err, "slice with nil values shouldn't have errors, just ignore value")
		require.Equal(t, doneTasks, nonNilTasks, "non-nil-tasks count should be equal to done-tasks")
	})

	t.Run("Run with empty slice", func(t *testing.T) {
		workersCount := 5
		maxErrorsCount := 1

		err := Run([]Task{}, workersCount, maxErrorsCount)

		require.NoError(t, err, "Run with empty slice shouldn't have any errors")
	})

	t.Run("Run with nil slice", func(t *testing.T) {
		workersCount := 5
		maxErrorsCount := 1

		err := Run(nil, workersCount, maxErrorsCount)

		require.NoError(t, err, "Run with nil slice shouldn't have any errors")
	})

	t.Run("Run with 0 or negative workers arg", func(t *testing.T) {
		tasks := []Task{}
		workersCount := 0
		maxErrorsCount := 1
		maxTasks := 30
		var doneTasks int64

		for i := 0; i < maxTasks; i++ {
			tasks = append(tasks, func() error {
				atomic.AddInt64(&doneTasks, 1)

				return nil
			})
		}

		err := Run(tasks, workersCount, maxErrorsCount)

		require.ErrorIs(t, err, ErrNoWorkers, "should return no-workers error if no workers found")
		require.Equal(t, int64(0), doneTasks, "tasks shouldn't be done without workers")
	})

	t.Run("Run with max 0 errors arg", func(t *testing.T) {
		tasks := []Task{}
		workersCount := 5
		maxErrorsCount := 0
		maxTasks := 30
		var doneTasks int64

		for i := 0; i < maxTasks; i++ {
			tasks = append(tasks, func() error {
				atomic.AddInt64(&doneTasks, 1)

				return nil
			})
		}

		err := Run(tasks, workersCount, maxErrorsCount)

		require.ErrorIs(t, err, ErrMaxErrorsIsNotValid, "should return is-not-valid-number error if max-errors number is too small")
		require.Equal(t, int64(0), doneTasks, "tasks shouldn't be done without workers")
	})
}
