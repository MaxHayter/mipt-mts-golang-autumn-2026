package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	before := runtime.NumGoroutine()
	stop := make(chan struct{})

	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-stop
		}()
	}

	close(stop) // Один сигнал завершает всех ожидающих workers.
	wg.Wait()
	after := runtime.NumGoroutine()
	fmt.Printf("goroutines: before=%d, after=%d\n", before, after)
}
