# Worker Pool
The library implements worker pool pattern by channels. Now realized 2 kinds of worker pool: `Pool` and `TimedPool`. `Pool` is the generic realization of worker pool pattern when tasks are received from outside. `TimedPool` - is the worker pool implementation where tasks are received by dispatcher on ticker event.

## Install 

```bash
go get github.com/dipdup.net/workerpool
```

## Examples

Usage of `Pool`

```go
package main

import (
	"context"
	"log"
	"time"
)

func main() {
	pool := NewPool(worker, 2)

	ctx, cancel := context.WithCancel(context.Background())
	pool.Start(ctx)

	dispatcher(ctx, pool)

	time.Sleep(time.Minute)

	cancel()

	if err := pool.Close(); err != nil {
		log.Panic(err)
	}
}

func worker(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Second)
			log.Printf("hello, %s", name)
			return
		}
	}
}

func dispatcher(ctx context.Context, pool *Pool[string]) {
	for _, name := range []string{"John", "Mark", "Peter", "Mike"} {
		select {
		case <-ctx.Done():
			return
		default:
			pool.AddTask(name)
		}
	}
}
```

Usage of `TimedPool`

```go
package main

import (
	"context"
	"log"
	"time"
)

func main() {
	pool := NewTimedPool(dispatcher, worker, nil, 2, 60*1000) // tasks will be received over 60 seconds

	ctx, cancel := context.WithCancel(context.Background())
	pool.Start(ctx)

	time.Sleep(time.Minute)

	cancel()

	if err := pool.Close(); err != nil {
		log.Panic(err)
	}
}

func worker(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Second)
			log.Printf("hello, %s", name)
			return
		}
	}
}

func dispatcher(ctx context.Context) ([]string, error) {
	return []string{"John", "Mark", "Peter", "Mike"}, nil
}
```

`NewTimedPool` receives as arguments 3 handlers: `Dispatcher`, `Worker` and `ErrorHandler`. 

```go
// Worker - worker handler. Calls on task arriving.
type Worker[Task any] func(ctx context.Context, task Task)

// Dispatcher - the function is called by timer for receiving tasks
type Dispatcher[Task any] func(ctx context.Context) ([]Task, error)

// ErrorHandler - the function is called when error occured in dispatcher
type ErrorHandler[Task any] func(ctx context.Context, err error)
```

Also it receives 2 integers: workers count and time between dispatcher calls in milliseconds.