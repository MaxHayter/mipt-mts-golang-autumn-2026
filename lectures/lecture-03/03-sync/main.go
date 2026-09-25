package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	mu sync.RWMutex
	n  int
}

func (c *Counter) Add(delta int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.n += delta
}

func (c *Counter) Value() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.n
}

func main() {
	var c Counter
	var wg sync.WaitGroup

	for range 50 {
		wg.Go(func() {
			for range 1000 {
				c.Add(1)
			}
		})
	}
	
	wg.Wait()
	fmt.Println("итог:", c.Value())
}
