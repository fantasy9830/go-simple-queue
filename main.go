package main

import (
	"context"
	"go-simple-queue/pkg/queue"
	"log"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue.Init(ctx)

	var wg sync.WaitGroup
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		queue.AddWork(queue.Job{
			TaskFunc: func(num int) func(c context.Context) error {
				return func(c context.Context) error {
					defer wg.Done()
					log.Println(num)

					return nil
				}
			}(i),
		})
	}

	wg.Wait()
}
