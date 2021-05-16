package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	fulfilled := make(chan struct{})
	out := make(Bi, 100000)
	outEmpty := make(Bi)
	close(outEmpty)

	if in == nil {
		close(out)

		return out
	}

	go func() {
		defer close(fulfilled)
		defer close(out)

		for _, stage := range stages {
			in = stage(in)
		}

		for item := range in {
			select {
			case <-done:
				break
			case out <- item:
			}
		}
	}()

	for {
		select {
		case <-done:
			return outEmpty
		case <-fulfilled:
			return out
		}
	}
}
