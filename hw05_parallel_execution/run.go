package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type ErS struct {
	mutex       sync.Mutex
	errorsCount int
	errorsLimit int
	stopTasks   bool
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, workersNum, errorsLimit int) error {
	taskChannel := make(chan Task, len(tasks))
	for _, task := range tasks {
		taskChannel <- task
	}
	close(taskChannel)
	// allDoneChannel := make(chan interface{})
	// errorChannel := make(chan error, workersNum)
	errorChannel := make(chan error)

	wg := sync.WaitGroup{}
	wg.Add(len(tasks))

	mutex := sync.Mutex{}

	errorsCount := 0
	stopTasks := false

	go func() {
		_, ok := <-errorChannel
		if ok {
			for {
				_, ok := <-taskChannel
				if !ok {
					return
				}

				wg.Done()
			}
		}
	}()

	for i := 0; i < workersNum; i++ {
		go func() {
			for {
				taskFunc, ok := <-taskChannel
				if !ok {
					return
				}
				taskError := taskFunc()
				wg.Done()

				mutex.Lock()
				if taskError != nil && errorsLimit > 0 {
					errorsCount++
					if errorsCount >= errorsLimit {
						if stopTasks {
							mutex.Unlock()
							return
						}
						stopTasks = true
						errorChannel <- ErrErrorsLimitExceeded
						close(errorChannel)
						fmt.Print("*")
						mutex.Unlock()
						return
					}
				}
				mutex.Unlock()
			}
		}()
	}

	wg.Wait()

	if stopTasks {
		return ErrErrorsLimitExceeded
	} else {
		errorChannel <- ErrErrorsLimitExceeded
	}

	return nil
}
