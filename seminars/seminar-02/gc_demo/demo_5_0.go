package main

import (
	"fmt"
	"runtime"
	"time"
)

func allocateMemory(iterations int) {
	for i := 0; i < iterations; i++ {
		// Выделяем память, которая станет мусором
		s := make([]int64, 10000)

		// Имитация работы: заполняем значениями
		for i := range s {
			s[i] = 0
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
