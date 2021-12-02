package queue

import (
	"context"
	"log"
)

// Worker represents the worker that executes the job
type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
	}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start(ctx context.Context) {
	go func(ctx context.Context) {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case <-ctx.Done():
				return
			case job := <-w.JobChannel:
				// we have received a work request.
				if err := job.TaskFunc(ctx); err != nil {
					log.Println(err.Error())
				}
			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}(ctx)
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
