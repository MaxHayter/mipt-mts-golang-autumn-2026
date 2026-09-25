package main

import (
	"fmt"
	"sync"
	"time"
)

func race0() {
	var n int

	for range 1000 {
		go func() {
			n++
		}()
	}

	time.Sleep(100 * time.Millisecond)

	fmt.Println(n)
}

func race1() {
	counters := map[int]int{}
	for i := range 5 {
		go func(counters map[int]int, th int) {
			for j := range 5 {
				counters[th*10+j]++
			}
		}(counters, i)
	}
	time.Sleep(100 * time.Millisecond)

	fmt.Println("counters result", counters)
}

func race2() {
	counters := map[int]int{}
	mu := &sync.Mutex{}
	for i := range 5 {
		go func(counters map[int]int, th int, mu *sync.Mutex) {
			for j := range 5 {
				mu.Lock()
				counters[th*10+j]++
				mu.Unlock()
			}
		}(counters, i, mu)
	}
	time.Sleep(100 * time.Millisecond)
	mu.Lock()
	fmt.Println("counters result", counters)
	mu.Unlock()
}

func main() {
	// race0()

	// race1()

	// race2()
}
