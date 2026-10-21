package main

import (
	"fmt"
	"sync"
)

func main() {
	var count int
	var wg sync.WaitGroup

	for range 1_000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count++ // WaitGroup не защищает общую память.
		}()
	}

	wg.Wait()
	fmt.Println("count:", count)
}
