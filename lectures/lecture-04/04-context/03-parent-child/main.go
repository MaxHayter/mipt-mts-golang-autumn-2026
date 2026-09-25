package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func worker(ctx context.Context, name string) {
	<-ctx.Done()
	fmt.Println(name, "остановлен:", ctx.Err())
}

// child1 и child2 порождены от одного parent. Отменяем только parent —
// оба ребёнка получают Done(), хотя cancel у них никто не вызывал.
func main() {
	parent, cancelParent := context.WithCancel(context.Background())
	child1, cancel1 := context.WithCancel(parent)
	child2, cancel2 := context.WithCancel(parent)
	defer cancel1()
	defer cancel2()

	var wg sync.WaitGroup
	wg.Go(func() { worker(child1, "child-1") })
	wg.Go(func() { worker(child2, "child-2") })

	time.Sleep(50 * time.Millisecond)
	fmt.Println("main: отменяем родительский контекст")
	cancelParent()
	wg.Wait()
}
