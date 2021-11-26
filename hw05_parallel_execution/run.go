package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Consumer(wg *sync.WaitGroup, ch <-chan Task, m int, errCount *int32) { // Функция-потребитель
	defer wg.Done()
	for {
		task, ok := <-ch
		if !ok { // Если канал закрыт, завершаем функцию
			return
		}
		if err := task(); err != nil && m > 0 { // Если m <= 0 игнорируем ошибки
			atomic.AddInt32(errCount, 1)
			return
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// Инициализация переменных
	var wg sync.WaitGroup // WG
	ch := make(chan Task) // Канал передачи данных
	var errResult error   // Конечная ошибка
	var errCount int32    // Счётчик ошибок, тип int64, так как используются атомарные операции

	// Если количество горутин больше длины массива, то количество горутин должно быть ограничено
	if n > len(tasks) {
		n = len(tasks)
	}

	wg.Add(n)
	for i := 0; i < n; i++ { // Создём горутины-потребители
		go Consumer(&wg, ch, m, &errCount)
	}

	// Производитель.
	// Передаём в канал функции из слайса.
	// Если количество ошибок больше количества горутин, прерываем цикл.
	// Если количество ошибок больше предельного количества ошибок, прерываем цикл.
	// Если предельное количество ошибок - 0 или меньше, ошибки игнорируются
	i := 0
Producer:
	for i < len(tasks) {
		select {
		case ch <- tasks[i]:
			i++
		default:
			if (int(atomic.LoadInt32(&errCount)) >= n || int(atomic.LoadInt32(&errCount)) >= m) && m > 0 {
				errResult = ErrErrorsLimitExceeded
				break Producer
			}
		}
	}
	close(ch)
	wg.Wait()
	return errResult
}
