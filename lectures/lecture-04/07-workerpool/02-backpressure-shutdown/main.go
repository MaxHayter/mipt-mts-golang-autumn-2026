package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type job struct {
	id int
}

func worker(ctx context.Context, id int, jobs <-chan job) {
	for {
		select {
		case j, ok := <-jobs:
			if !ok {
				return
			}
			time.Sleep(20 * time.Millisecond) // имитация чтения файла с диска
			fmt.Printf("воркер %d обработал задачу %d\n", id, j.id)
		case <-ctx.Done():
			fmt.Printf("воркер %d остановлен: %v\n", id, ctx.Err())
			return
		}
	}
}

func main() {
	const numWorkers = 3
	const queueSize = 2 // маленький буфер — это и есть backpressure

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	jobs := make(chan job, queueSize)
	var wg sync.WaitGroup
	for i := 1; i <= numWorkers; i++ {
		wg.Go(func() { worker(ctx, i, jobs) })
	}

	go func() {
		for i := 1; i <= 20; i++ {
			select {
			case jobs <- job{id: i}: // блокируется, если у воркеров нет места
			case <-ctx.Done():
				return
			}
		}
		close(jobs)
	}()

	wg.Wait()
	fmt.Println("пул остановлен")
}
