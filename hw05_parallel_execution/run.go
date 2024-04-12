package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, workersNum, errorsLimit int) error {
	taskChannel := make(chan Task, len(tasks))
	for _, task := range tasks {
		taskChannel <- task
	}
	close(taskChannel)
	errorChannel := make(chan error)

	tasksDoneWaiter := sync.WaitGroup{}
	tasksDoneWaiter.Add(len(tasks))

	closeTasksWaiter := sync.WaitGroup{}
	closeTasksWaiter.Add(1)

	mutex := sync.Mutex{}

	errorsCount := 0
	stopTasks := false

	go EndAllTasks(taskChannel, errorChannel, &closeTasksWaiter, &tasksDoneWaiter)

	for i := 0; i < workersNum; i++ {
		go TaskProcessing(taskChannel, errorChannel, errorsLimit, &mutex, &tasksDoneWaiter, &stopTasks, &errorsCount)
	}

	tasksDoneWaiter.Wait()

	if stopTasks {
		return ErrErrorsLimitExceeded
	}

	errorChannel <- ErrErrorsLimitExceeded
	closeTasksWaiter.Wait()

	return nil
}

func EndAllTasks(taskChannel chan Task, errorChannel chan error, closeTasksWaiter, tasksDoneWaiter *sync.WaitGroup) {
	_, ok := <-errorChannel
	if ok {
		for {
			_, ok := <-taskChannel
			if !ok {
				closeTasksWaiter.Done()
				return
			}

			tasksDoneWaiter.Done()
		}
	}

	closeTasksWaiter.Done()
}

func TaskProcessing(taskChannel chan Task, errorChannel chan error, errorsLimit int, mutex *sync.Mutex,
	tasksDoneWaiter *sync.WaitGroup, stopTasks *bool, errorsCount *int) {
	for {
		taskFunc, ok := <-taskChannel
		if !ok {
			return
		}
		taskError := taskFunc()
		tasksDoneWaiter.Done()

		mutex.Lock()
		if taskError != nil && errorsLimit > 0 {
			*errorsCount++
			if *errorsCount >= errorsLimit {
				if *stopTasks {
					mutex.Unlock()
					return
				}
				*stopTasks = true
				errorChannel <- ErrErrorsLimitExceeded
				close(errorChannel)
				mutex.Unlock()
				return
			}
		}
		mutex.Unlock()
	}
}
