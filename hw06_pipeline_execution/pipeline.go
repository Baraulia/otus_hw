package hw06pipelineexecution

import (
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		out := make(Bi)
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func(in In, out Bi, done In) {
			defer close(out)
			wg.Done()
			for value := range in {
				select {
				case <-done:
					return
				default:
					out <- value
				}
			}
		}(in, out, done)

		wg.Wait()
		//time.Sleep(time.Millisecond * 1)
		in = stage(out)
	}

	return in
}
