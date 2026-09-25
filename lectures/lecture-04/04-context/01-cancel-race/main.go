package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

// Десять студентов одновременно думают над задачей. Как только один
// сдаёт первым, остальным больше нет смысла думать — общий cancel()
// останавливает всех разом.
func student(ctx context.Context, id int, done chan<- int) {
	thinking := time.Duration(rand.IntN(100)+10) * time.Millisecond
	select {
	case <-ctx.Done():
		fmt.Printf("студент %d прекратил думать: %v\n", id, ctx.Err())
	case <-time.After(thinking):
		select {
		case done <- id:
		case <-ctx.Done():
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan int, 1)

	var wg sync.WaitGroup
	for id := range 10 {
		wg.Go(func() { student(ctx, id, done) })
	}

	first := <-done
	fmt.Println("первым сдал студент", first)
	cancel()
	wg.Wait()
}
