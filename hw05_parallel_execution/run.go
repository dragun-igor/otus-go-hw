package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Consumer(wg *sync.WaitGroup, ch <-chan Task, done <-chan bool, m int, errCount *int32) { // Функция-приёмник
	defer wg.Done()
	for {
		select {
		case <-done:
			return
		case task := <-ch:
			if err := task(); err != nil && m > 0 { // Если m <= 0 игнорируем ошибки
				atomic.AddInt32(errCount, 1)
				return
			}
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	ch := make(chan Task)   // Канал передачи функций
	done := make(chan bool) // Канал для передачи команды завершения горутины
	var workerAmount int    // Количество воркеров
	var errResult error     // Конечная ошибка
	var errCount int32      // Счётчик ошибок
	if n > len(tasks) {     // Определяем количество воркеров
		workerAmount = len(tasks)
	} else {
		workerAmount = n
	}
	for i := 0; i < workerAmount; i++ { // Создём горутины
		wg.Add(1)
		go Consumer(&wg, ch, done, m, &errCount)
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
			if int(atomic.LoadInt32(&errCount)) >= workerAmount && m > 0 {
				errResult = ErrErrorsLimitExceeded
				break Producer
			}
		}
	}
	close(done)
	wg.Wait()
	close(ch)
	return errResult
}
