package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// outChan is function for processing the output channel.
func outChan(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer func() {
			close(out)
			for range in {
			}
		}()
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

// ExecutePipeline executes stages.
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		out = stage(outChan(out, done))
	}
	return outChan(out, done)
}
