package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

var wg sync.WaitGroup

func Consumer(ch chan Task, done chan bool, m int, errCount *int32) { // Функция-приёмник
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
	ch := make(chan Task)   // Канал передачи функций
	done := make(chan bool) // Канал для передачи команды завершения горутины
	var workerAmount int    // Количество воркеров
	var errResult error     // Конечная ошибка
	var errCount int32

	if n > len(tasks) { // Определяем количество воркеров
		workerAmount = len(tasks)
	} else {
		workerAmount = n
	}
	wg.Add(workerAmount)                // Задаём количество горутин для вэйтгруппы
	for i := 0; i < workerAmount; i++ { // Создём горутины
		go Consumer(ch, done, m, &errCount)
	}
	i := 0
Producer:
	for { // Мне очень не нравится этот кусок, но по-другому не придумал
		select {
		case ch <- tasks[i]:
			if i < len(tasks)-1 {
				i++
			} else {
				break Producer
			}
		default:
			if int(errCount) >= workerAmount && m > 0 {
				errResult = ErrErrorsLimitExceeded
				break Producer
			}
		}
	}

	for i := 0; i < workerAmount-int(errCount); i++ { // Раздём команды на завершение работы горутин
		done <- true
	}
	wg.Wait()
	close(ch)
	close(done)
	return errResult
}
