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

	closeStream := make(chan bool)

	stopWrite := false

	mutex := sync.Mutex{}

	go func() {
		for s := range out {
			mutex.Lock()
			if !stopWrite {
				valueStream <- s
				mutex.Unlock()
				continue
			}
			mutex.Unlock()
			break
		}

		closeStream <- true
	}()

	go func() {
		<-done
		closeStream <- true
	}()

	go func() {
		<-closeStream
		mutex.Lock()
		stopWrite = true
		close(valueStream)
		mutex.Unlock()
	}()

	return valueStream
}
