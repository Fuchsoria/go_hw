package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)
	defer close(out)

	if in == nil {
		return out
	}

	for _, stage := range stages {
		if done != nil {
			<-done

			return out
		}

		in = stage(in)
	}

	return in
}
