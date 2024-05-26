package hw06pipelineexecution

import (
	"fmt"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}

	out := doStage(done, stages[0], in)
	for i := 1; i < len(stages); i++ {
		out = doStage(done, stages[i], out)
	}

	return out
}

func doStage(done In, stage Stage, in In) Out {
	valueStream := make(Bi)
	out := stage(in)

	go func() {
		for s := range out {
			select {
			case <-done:
				close(valueStream)
				return
			default:
				fmt.Printf("stage: %v\n", s)
				valueStream <- s
			}
		}

		close(valueStream)
	}()

	return valueStream
}
