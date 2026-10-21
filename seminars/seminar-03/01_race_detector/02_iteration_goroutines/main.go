package main

import (
	"fmt"
	"time"
)

func main() {
	var count int

	for range 1_000 {
		go func() {
			count++ // Одновременная запись в общую память.
		}()
	}

	// Это не синхронизация, а случайная пауза.
	time.Sleep(10 * time.Millisecond)
	fmt.Println("count:", count)
}
