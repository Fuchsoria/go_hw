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
			in = stage(in)
		}

		for item := range in {
			out <- item
		}

		close(out)
	}()

	return out
}
