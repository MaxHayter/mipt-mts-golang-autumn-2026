package main

import (
	"fmt"
	"sync"
	"time"
)

// close(ready) будит разом всех, кто заблокирован на приёме из ready —
// не по одному, а одновременным пробуждением всей очереди ожидающих
// горутин внутри рантайма. Ровно на этом механизме позже построен
// context.Done().
func main() {
	const workers = 5
	ready := make(chan struct{})

	var wg sync.WaitGroup
	unblockedAt := make([]time.Duration, workers)

	start := time.Now()
	for i := range workers {
		wg.Go(func() {
			<-ready
			elapsed := time.Since(start)
			unblockedAt[i] = elapsed
		})
	}

	time.Sleep(50 * time.Millisecond) // дать всем воркерам дойти до <-ready
	close(ready)                      // будим все пять горутин одним close
	wg.Wait()

	for id, d := range unblockedAt {
		fmt.Printf("воркер %d проснулся через %v после close\n", id, d)
	}
}
