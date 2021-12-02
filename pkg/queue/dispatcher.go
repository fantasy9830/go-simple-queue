package queue

import (
	"context"
)

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	workerpool chan chan Job

	maxWorkers int
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)

	return &Dispatcher{workerpool: pool, maxWorkers: maxWorkers}
}

func (d *Dispatcher) Run(ctx context.Context) {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.workerpool)
		worker.Start(ctx)
	}

	go d.dispatch(ctx)
}

func (d *Dispatcher) dispatch(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-JobQueue:
			if !ok {
				break
			}

			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.workerpool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}
