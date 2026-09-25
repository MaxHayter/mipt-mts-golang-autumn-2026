package main

import (
	"fmt"
	"sync"
)

// Самый голый воркер-пул: N горутин читают из одного канала-очереди,
// пока его не закроют. Никакого контекста, никакого backpressure —
// только механика «общий канал, фиксированное число читателей».
func worker(id int, jobs <-chan int) {
	for job := range jobs {
		fmt.Printf("воркер %d обработал задачу %d\n", id, job)
	}
}

func main() {
	const numWorkers = 3
	jobs := make(chan int)

	var wg sync.WaitGroup
	for i := 1; i <= numWorkers; i++ {
		wg.Go(func() { worker(i, jobs) })
	}

	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
	fmt.Println("все задачи разобраны")
}
