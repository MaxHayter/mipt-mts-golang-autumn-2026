package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	const goroutines = 100
	const perGoroutine = 10000

	// var counter atomic.Int64
	var counter int64
	var wg sync.WaitGroup

	for range goroutines {
		wg.Go(func() {
			for range perGoroutine {
				// counter.Add(1)
				atomic.AddInt64(&counter, 1)
			}
		})
	}
	wg.Wait()

	fmt.Println("итог:", counter /*.Load()*/, "ожидалось:", goroutines*perGoroutine)

	counter = 0
	for i := range 10 {
		wg.Go(func() {
			swapped := atomic.CompareAndSwapInt64(&counter, 0, 100)
			fmt.Printf("i = %d, swapped = %t\n", i, swapped)
		})
	}

	wg.Wait()
}
