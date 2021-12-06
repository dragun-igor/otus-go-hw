package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

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
					temp <- val
				}
			}
		}(temp, out)
		out = stage(temp)
	}
	return out
}
