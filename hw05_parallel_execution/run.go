package hw05parallelexecution

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, goroutinesNum, errorsLimit int) error {
	ch := make(chan Task, len(tasks))
	for _, task := range tasks {
		ch <- task
	}
	close(ch)

	wg := sync.WaitGroup{}
	wg.Add(len(tasks))

	mutex := sync.Mutex{}

	errorsCount := 0
	stopTasks := false

	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < goroutinesNum; i++ {
		go func(ctx context.Context) {
			for {
				// select {
				// case <-ctx.Done():
				// 	wg.Done()
				// 	return
				// case taskFunc := <-ch:
				// 	fmt.Println("Task Start")
				// 	taskError := taskFunc()
				// 	fmt.Println("Task End")
				// 	if taskError != nil && errorsLimit > 0 {
				// 		mutex.Lock()
				// 		errorsCount++
				// 		if errorsCount >= errorsLimit {
				// 			cancel()
				// 		}
				// 		mutex.Unlock()
				// 	}
				// 	wg.Done()
				// }

				taskFunc, ok := <-ch
				if !ok {
					fmt.Print("?")
					return
				}

				if stopTasks {
					wg.Done()
					fmt.Print("s")
					continue
				}

				taskError := taskFunc()
				if taskError != nil && errorsLimit > 0 {
					mutex.Lock()
					errorsCount++
					if errorsCount >= errorsLimit {
						// cancel()
						fmt.Print("!")
						stopTasks = true
					}
					mutex.Unlock()
					fmt.Print(errorsCount)
				}

				fmt.Print("+")
				wg.Done()
			}

		}(ctx)
	}

	wg.Wait()
	cancel()

	if stopTasks {
		return ErrErrorsLimitExceeded
	}

	return nil
}
