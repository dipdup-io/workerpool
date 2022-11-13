package workerpool

import (
	"context"
	"time"
)

// Dispatcher - the function is called by timer for receiving tasks
type Dispatcher[Task any] func(ctx context.Context) ([]Task, error)

// ErrorHandler - the function is called when error occured in dispatcher
type ErrorHandler[Task any] func(ctx context.Context, err error)

// TimedPool - worker pool which receives task by executing `dispatcher` function by ticker.
type TimedPool[Task any] struct {
	*Pool[Task]

	dispatcher           Dispatcher[Task]
	errorHandler         ErrorHandler[Task]
	periodReceivingTasks time.Duration
}

// NewTimedPool - creates new `TimedPool`.
// `dispatcher` - is a function returning tasks. If it's null `Start` function immidiately exits.
// `worker` - is a job executing by a task.
// `errorHandler` - callback to error handling in `dispatcher` executing. May be null.
// `workersCount` - count of workers executing jobs.
// `periodReceivingTasks` - ticker period in milliseconds when `dispatcher` will be called.
func NewTimedPool[Task any](
	dispatcher Dispatcher[Task],
	worker Worker[Task],
	errorHandler ErrorHandler[Task],
	workersCount int,
	periodReceivingTasks int,

) *TimedPool[Task] {
	return &TimedPool[Task]{
		Pool:                 NewPool(worker, workersCount),
		dispatcher:           dispatcher,
		errorHandler:         errorHandler,
		periodReceivingTasks: time.Millisecond * time.Duration(periodReceivingTasks),
	}
}

// Start - starts pool and dispatcher.
func (pool *TimedPool[Task]) Start(ctx context.Context) {
	if pool.dispatcher == nil {
		return
	}

	pool.Pool.Start(ctx)

	pool.wg.Add(1)
	go pool.dispatch(ctx)
}

func (pool *TimedPool[Task]) dispatch(ctx context.Context) {
	defer pool.wg.Done()

	ticker := time.NewTicker(pool.periodReceivingTasks)
	defer ticker.Stop()

	// First tick
	pool.tick(ctx, ticker)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pool.tick(ctx, ticker)
		}
	}
}

func (pool *TimedPool[Task]) tick(ctx context.Context, ticker *time.Ticker) {
	tasks, err := pool.dispatcher(ctx)
	if err != nil && pool.errorHandler != nil {
		pool.errorHandler(ctx, err)
		return
	}

	if len(tasks) == 0 {
		ticker.Reset(pool.periodReceivingTasks * 10)
		return
	}

	for i := range tasks {
		select {
		case <-ctx.Done():
			return
		default:
			pool.Pool.AddTask(tasks[i])
		}
	}
	ticker.Reset(pool.periodReceivingTasks)
}
