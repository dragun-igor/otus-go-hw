package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func selectCase() {

}

func outChan(in In, done In) Bi {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
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
	}()
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		temp := make(Bi)
		go func(temp Bi, out Out) {
			defer close(temp)
			for {
				select {
				case <-done:
					return
				case val, ok := <-out:
					if !ok {
						return
					}
					select {
					case <-done:
						return
					case temp <- val:
					}
				}
			}
		}(temp, out)
		out = stage(temp)
	}

	// Добавляем выходной канал на чтение и запись, чтобы по сигналу done его можно было закрыть и прервать пайплайн
	return outChan(out, done)
}
