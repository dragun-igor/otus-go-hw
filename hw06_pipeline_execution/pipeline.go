package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	stageAmount := len(stages)
	results := make([]Bi, 0, stageAmount+1)
	for i := 0; i < len(stages)+1; i++ {
		results = append(results, make(Bi))
	}

	go func() {
		<-done
	}()

	// Переброска значений в первую ячейку слайса, после окончания закрываем канал
	go func() {
		defer close(results[0])
		for {
			select {
			case <-done:
				return
				case val, ok := <-in:
					if !ok {
						return
					}
					results[0] <- val
			}
		}
	}()

	for i, stage := range stages {
		out := stage(results[i])
		go func(i int) {
			defer close(results[i+1])
			for {
				select {
				case <-done:
					return
				case val, ok := <-out:
					if !ok {
						return
					}
					results[i+1] <- val
				}
			}
		}(i)
	}

	// Возвращаем канал из последней ячейки слайса
	return results[stageAmount]
}
