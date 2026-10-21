package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	before := runtime.NumGoroutine()
	neverSent := make(chan struct{})

	for range 100 {
		go func() {
			<-neverSent // Ждём событие, которое никогда не наступит.
		}()
	}

	time.Sleep(50 * time.Millisecond)
	after := runtime.NumGoroutine()
	fmt.Printf("goroutines: before=%d, after=%d\n", before, after)
}
