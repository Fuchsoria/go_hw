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

	t.Run("mixed random tasks", func(t *testing.T) {
		tasks := []Task{}
		maxTasks := 100
		workersCount := 10
		maxErrorsCount := 50
		var doneTasks int64
		var totalDoneTasks int64
		var errorTasks int64
		var totalErrorTasks int64

		randBool := func() bool {
			rand.Seed(time.Now().UnixNano())
			return rand.Intn(2) == 1
		}

		for i := 0; i < maxTasks; i++ {
			if shouldDone := randBool(); shouldDone {
				atomic.AddInt64(&totalDoneTasks, 1)

				tasks = append(tasks, func() error {
					atomic.AddInt64(&doneTasks, 1)

					return nil
				})
			} else {
				atomic.AddInt64(&totalErrorTasks, 1)

				tasks = append(tasks, func() error {
					atomic.AddInt64(&errorTasks, 1)

					return fmt.Errorf("error from task %d", i)
				})
			}
		}

		err := Run(tasks, workersCount, maxErrorsCount)

		if totalErrorTasks >= int64(maxErrorsCount) {
			fmt.Println(maxErrorsCount, totalErrorTasks, errorTasks, err)
			require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
			require.True(t, int64(maxErrorsCount) <= errorTasks, "not all errors tasks are done")
			require.True(t, errorTasks <= int64(maxErrorsCount+workersCount), "not all errors tasks are done")
		} else {
			require.NoError(t, err)
			require.Equal(t, totalDoneTasks, doneTasks, "not all positive tasks are done")
		}
	})

	t.Run("Run with nil values", func(t *testing.T) {
		tasks := []Task{}
		maxTasks := 50
		workersCount := 10
		maxErrorsCount := 1
		var doneTasks int64
		var totalNilTasks int64

		randBool := func() bool {
			rand.Seed(time.Now().UnixNano())
			return rand.Intn(2) == 1
		}

		for i := 0; i < maxTasks; i++ {
			if shouldDone := randBool(); shouldDone {
				tasks = append(tasks, func() error {
					atomic.AddInt64(&doneTasks, 1)

					return nil
				})
			} else {
				atomic.AddInt64(&totalNilTasks, 1)

				tasks = append(tasks, nil)
			}
		}

		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err, "slice with nil values shouldn't have errors, just ignore value")
		require.Equal(t, doneTasks+totalNilTasks, int64(maxTasks), "total tasks + nils should be equal to max tasks")
	})
}
