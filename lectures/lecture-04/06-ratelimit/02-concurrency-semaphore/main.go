package main

import (
	"fmt"
	"sync"
	"time"
)

// Ограничиваем не скорость по времени, а число одновременно
// работающих горутин: буферизованный канал как семафор — слот занят,
// пока в нём лежит значение. Ровно на этом устроен errgroup.SetLimit.
func main() {
	const maxConcurrent = 2
	sem := make(chan struct{}, maxConcurrent)

	var wg sync.WaitGroup
	for id := range 5 {
		wg.Go(func() {
			sem <- struct{}{}
			defer func() { <-sem }()

			fmt.Printf("задача %d начала работу\n", id)
			fmt.Printf("задача %d закончила\n", id)
			time.Sleep(1500 * time.Millisecond)
		})
	}
	wg.Wait()
}
