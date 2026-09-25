package main

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	iterationsNum = 6
	goroutinesNum = 6
)

func doWork(th int) {
	for j := range iterationsNum {
		// if j%2 != 0 {
		// 	time.Sleep(1 * time.Millisecond)
		// }
		fmt.Print(formatWork(th, j))
		// runtime.Gosched()
	}
}

func formatWork(in, j int) string {
	return fmt.Sprintln(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNum-in),
		"th", in,
		"iter", j, strings.Repeat("■", j))
}

func main() {
	// runtime.GOMAXPROCS(1)
	// for i := range goroutinesNum {
	// 	go doWork(i)
	// }

	// time.Sleep(100 * time.Millisecond)

	const n = 100_000
	release := make(chan struct{})
	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)
	for range n {
		go func() { <-release }()
	}
	runtime.ReadMemStats(&after)
	fmt.Printf("куча на горутину: ~%d байт\n", (after.HeapAlloc-before.HeapAlloc)/n)
	fmt.Printf("стек на горутину: ~%d байт\n", (after.StackInuse-before.StackInuse)/n)
	close(release) // отпускаем все n горутин разом, программа может завершатьс
}
