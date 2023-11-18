package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil {
		out := make(Bi)
		close(out)
		return out
	}

	for _, stage := range stages {
		out := make(Bi)
		go func(in In, out Bi, done In) {
			defer close(out)
			for {
				select {
				case <-done:
					return
				case value, ok := <-in:
					if !ok {
						return
					}
					out <- value
				}
			}
		}(in, out, done)

		in = stage(out)
	}

	return in
}
