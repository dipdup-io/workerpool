package workerpool

import (
	"context"
	"sync"
)

// Worker - worker handler. Calls on task arriving.
type Worker[Task any] func(ctx context.Context, task Task)

// Pool - worker pool entity
type Pool[Task any] struct {
	worker Worker[Task]
	tasks  chan Task
	wg     *sync.WaitGroup
}

// NewPool - creates new worker pool. `worker` - is a job. `workersCount` - count of workers executing jobs.
func NewPool[Task any](worker Worker[Task], workersCount int) *Pool[Task] {
	return &Pool[Task]{
		worker: worker,
		tasks:  make(chan Task, workersCount),
		wg:     new(sync.WaitGroup),
	}
}

// WorkersCount - returns workers count in the pool
func (pool *Pool[Task]) WorkersCount() int {
	return cap(pool.tasks)
}

// Close - wait returnning from workers and closing pool object
func (pool *Pool[Task]) Close() error {
	pool.wg.Wait()
	close(pool.tasks)
	return nil
}

// Start - begin worker pool
func (pool *Pool[Task]) Start(ctx context.Context) {
	if pool.worker == nil {
		return
	}

	for i := 0; i < pool.WorkersCount(); i++ {
		pool.wg.Add(1)
		go pool.start(ctx)
	}
}

func (pool *Pool[Task]) start(ctx context.Context) {
	defer pool.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case task := <-pool.tasks:
			pool.worker(ctx, task)
		}
	}
}

// AddTask - adds task to worker pool. Blocks if worker pool is full.
func (pool *Pool[Task]) AddTask(task Task) {
	pool.tasks <- task
}

// QueueSize - current count of tasks in pool
func (pool *Pool[Task]) QueueSize() int {
	return len(pool.tasks)
}
