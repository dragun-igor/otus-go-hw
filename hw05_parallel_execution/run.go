package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

var wg sync.WaitGroup

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	ch := make(chan Task)
	done := make(chan bool)
	var errCount int
	var workerAmount int
	var errResult error

	if n > len(tasks) {
		workerAmount = len(tasks)
	} else {
		workerAmount = n
	}
	wg.Add(workerAmount)
	for i := 0; i < workerAmount; i++ {
		go func() {
			defer wg.Done()
		Consumer:
			for {
				select {
				case <-done:
					break Consumer
				case task := <-ch:
					if err := task(); err != nil && m > 0 {
						errCount++
						break Consumer
					}
				}
			}
		}()
	}

	i := 0
Producer:
	for { // Мне очень не нравится этот кусок
		select {
		case ch <- tasks[i]:
			if i < len(tasks)-1 {
				i++
			} else {
				break Producer
			}
		default:
			if errCount >= workerAmount && m > 0 {
				errResult = ErrErrorsLimitExceeded
				break Producer
			}
		}
	}

	for i := 0; i < workerAmount-errCount; i++ {
		done <- true
	}
	wg.Wait()
	close(ch)
	close(done)
	return errResult
}
