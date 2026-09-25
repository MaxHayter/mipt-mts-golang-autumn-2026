package main

import (
	"fmt"
	"time"
)

func chan1() {
	ch1 := make(chan int)

	go func(in chan int) {
		fmt.Println("GO: before read from chan")
		val := <-in
		fmt.Println("GO: get from chan", val)
		fmt.Println("GO: after read from chan")
	}(ch1)
	fmt.Println("MAIN: before put to chan")

	ch1 <- 42
	// ch1 <- 100500

	fmt.Println("MAIN: after put to chan")
	time.Sleep(100 * time.Millisecond)
}

func chan2() {
	in := make(chan int)

	go func(out chan<- int) {
		for i := 0; i <= 10; i++ {
			fmt.Println("before", i)
			out <- i
			fmt.Println("after", i)
		}

		close(out)
		fmt.Println("generator finish")
	}(in)

	time.Sleep(time.Second)

	// for i := range in {
	// 	fmt.Println("\tget", i)
	// }

	for {
		v, ok := <-in
		if !ok {
			break
		}
		fmt.Println("\tget", v)
	}

	time.Sleep(time.Second)
}

func main() {
	// chan1()

	chan2()
}
