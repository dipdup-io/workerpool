package workerpool

import (
	"context"
	"sync"
)

// Runs functions with one wait group
type Group struct {
	wg *sync.WaitGroup
}

// NewGroup - creates new Group
func NewGroup() Group {
	return Group{
		wg: new(sync.WaitGroup),
	}
}

// Go - runs function with wait group
func (g Group) Go(f func()) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		f()
	}()
}

// GoCtx - runs function with wait group using context
func (g Group) GoCtx(ctx context.Context, f func(ctx context.Context)) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		f(ctx)
	}()
}

// Wait - waits until grouped functions end
func (g Group) Wait() {
	g.wg.Wait()
}
