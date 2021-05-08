package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		if in == nil {
			return
		}

		for _, stage := range stages {
			in = stage(in)

			if done == nil {
				continue
			}

			if _, ok := <-done; !ok {
				return
			}
		}

		for item := range in {
			select {
			case <-done:
				break
			case out <- item:
			}
		}
	}()

	return out
}
