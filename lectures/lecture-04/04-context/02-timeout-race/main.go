package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

const budget = 50 * time.Millisecond

// Та же гонка, но вместо ручного cancel — общий дедлайн: кто не успел
// за budget, просто не отвечает, никто никого не отменяет руками.
func student(ctx context.Context, id int, done chan<- int) {
	thinking := time.Duration(rand.IntN(100)+10) * time.Millisecond
	select {
	case <-ctx.Done():
	case <-time.After(thinking):
		select {
		case done <- id:
		case <-ctx.Done():
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), budget)
	defer cancel()
	done := make(chan int, 10)

	var wg sync.WaitGroup
	for id := range 10 {
		wg.Go(func() { student(ctx, id, done) })
	}
	wg.Wait()
	close(done)

	finished := 0
	for range done {
		finished++
	}
	fmt.Printf("успели за %v: %d из 10 студентов\n", budget, finished)
}
