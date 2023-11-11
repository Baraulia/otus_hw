package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	outWithDone := make(Bi)

	go func() {
		defer close(outWithDone)
		for value := range in {
			select {
			case <-done:
				return
			default:
				outWithDone <- value
			}
		}
	}()

	return buildPipeline(outWithDone, done, stages...)
}

func buildPipeline(in Out, done In, stages ...Stage) Out {
	for _, stage := range stages {
		out := make(Bi)
		go func(in In, out Bi, done In) {
			defer close(out)
			for value := range in {
				select {
				case <-done:
					return
				default:
				}

				out <- value
			}
		}(in, out, done)

		in = stage(out)
	}

	return in
}
