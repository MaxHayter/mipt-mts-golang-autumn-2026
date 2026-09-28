package main

import (
	"fmt"
	"runtime"
	"time"
)

var bigSlice = make([]int64, 10000)

func allocateMemory(iterations int) {
	for i := 0; i < iterations; i++ {
		// Переиспользуем один и тот же слайс вместо создания нового
		for j := range bigSlice {
			bigSlice[j] = 0
		}

		if i%100000 == 0 {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			fmt.Printf("Alloc: %v MB, TotalAlloc: %v MB\n",
				m.Alloc/1024/1024, m.TotalAlloc/1024/1024)
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func main() {
	fmt.Println("Starting allocation test...")
	allocateMemory(1000000)
	fmt.Println("Done")
}
