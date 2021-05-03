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
		for _, stage := range stages {
			select {
			case <-done:
				return
			default:
			}

			in = stage(in)
		}

	L:
		for item := range in {
			select {
			case <-done:
				break L
			default:
			}

			select {
			case out <- item:
			}
		}

		close(out)
	}()

	return out
}
