package workerpool

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestGroupWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	group := NewGroup()
	group.GoCtx(ctx, ticker1)
	group.GoCtx(ctx, ticker2)

	time.Sleep(10 * time.Second)
	cancel()

	group.Wait()
}

func ticker1(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			log.Print("second")
		}
	}
}

func ticker2(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			log.Print("2 second")
		}
	}
}

func TestGroup(t *testing.T) {
	group := NewGroup()
	group.Go(ticker3)
	group.Go(ticker4)
	group.Wait()
}

func ticker3() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var count int

	for range ticker.C {
		count++
		log.Print("second")

		if count == 5 {
			return
		}
	}
}

func ticker4() {
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	var count int

	for range ticker.C {
		count++
		log.Print("2 second")

		if count == 5 {
			return
		}
	}
}
