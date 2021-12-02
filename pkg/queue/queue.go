package queue

import (
	"context"
)

var (
	MaxQueue  = 10000
	MaxWorker = 1024
)

// TaskFunc TaskFunc
type TaskFunc func(ctx context.Context) error

// Job represents the job to be run
type Job struct {
	TaskFunc TaskFunc
}

// A buffered channel that we can send work requests on.
var JobQueue chan Job

func Init(ctx context.Context) error {
	JobQueue = make(chan Job, MaxQueue)

	dispatcher := NewDispatcher(MaxWorker)
	dispatcher.Run(ctx)

	return nil
}

// AddWork Push the work onto the queue.
func AddWork(work Job) {
	JobQueue <- work
}
