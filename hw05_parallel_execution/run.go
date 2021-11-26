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
	// Инициализация переменных
	var wg sync.WaitGroup   // WG
	ch := make(chan Task)   // Канал передачи данных
	done := make(chan bool) // Сигнальный канал для завершения работы горутин
	var errResult error     // Конечная ошибка
	var errCount int32      // Счётчик ошибок

	if n > len(tasks) {
		n = len(tasks)
	}

	wg.Add(n)
	for i := 0; i < n; i++ { // Создём горутины
		go Consumer(&wg, ch, done, m, &errCount)
	}

	i := 0
Producer:
	for i < len(tasks) {
		select {
		case ch <- tasks[i]:
			i++
		default:
			if int(atomic.LoadInt32(&errCount)) >= n && m > 0 {
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
