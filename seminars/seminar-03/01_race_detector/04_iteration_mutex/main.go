package main

import (
	"fmt"
	"sync"
)

func main() {
	var (
		count int
		mu    sync.Mutex
		wg    sync.WaitGroup
	)

	for range 1_000 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			mu.Lock()
			count++
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println("count:", count)
}
