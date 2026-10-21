package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
)

var cache = make([][]byte, 0, 30)

func main() {
	for range 30 {
		buffer := make([]byte, 1<<20) // 1 MiB
		cache = append(cache, buffer[:1])
	}

	runtime.GC()
	printHeap("after keeping 30 one-byte subslices")
	writeHeapProfile("memory-initial.pprof")
}

func printHeap(label string) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	fmt.Printf("%s: HeapAlloc = %.1f MiB\n", label, float64(stats.HeapAlloc)/(1<<20))
}

func writeHeapProfile(path string) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := pprof.WriteHeapProfile(file); err != nil {
		panic(err)
	}
}
