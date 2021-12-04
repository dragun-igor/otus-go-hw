package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Consumer(wg *sync.WaitGroup, ch <-chan Task, errCount *int32) {
	defer wg.Done()
	for task := range ch {
		if err := task(); err != nil {
			atomic.AddInt32(errCount, 1)
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	ch := make(chan Task)
	var errResult error
	var errCount int32

	if n > len(tasks) {
		n = len(tasks)
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go Consumer(&wg, ch, &errCount)
	}

Producer:
	for i := 0; i < len(tasks); {
		select {
		case ch <- tasks[i]:
			i++
		default:
			if int(atomic.LoadInt32(&errCount)) >= m && m > 0 {
				break Producer
			}
		}
	}
	close(ch)
	wg.Wait()
	if int(errCount) >= m && m > 0 {
		errResult = ErrErrorsLimitExceeded
	}
	return errResult
}
