package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)

	if in == nil {
		close(out)

		return out
	}

	for _, stage := range stages {
		in = stage(in)
	}

	go func() {
		defer close(out)

		for {
			select {
			case <-done:
				return
			case item, ok := <-in:
				if !ok {
					return
				}

				select {
				case <-done:
					return
				case out <- item:
				}
			}
		}
	}()

	return out
}
