package hw06pipelineexecution

import "sync/atomic"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

var numGoroutines int64 // Переменная для теста

func outChan(in In, done In) Out {
	// Переброска значений в канал доступный для записи.
	out := make(Bi)
	go func() {
		defer func() {
			close(out)
			<-in
			atomic.AddInt64(&numGoroutines, -1)
		}()
		atomic.AddInt64(&numGoroutines, 1)
		for {
			select {
			case <-done:
				return
			default:
				select {
				case <-done:
					return
				case val, ok := <-in:
					if !ok {
						return
					}
					select {
					case <-done:
						return
					case out <- val:
					}
				}
			}
		}
	}()
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Выполнение пайплайна
	out := in
	for _, stage := range stages {
		out = stage(outChan(out, done))
	}
	return outChan(out, done)
}
