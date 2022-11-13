# Worker Pool
The library implements worker pool pattern by channels. Now realized 2 kinds of worker pool: `Pool` and `TimedPool`. `Pool` is the generic realization of worker pool pattern when tasks are received from outside. `TimedPool` - is the worker pool iplementation where tasks are received by dispatcher on ticker event.

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
