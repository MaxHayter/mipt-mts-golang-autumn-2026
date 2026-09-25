package main

import (
	"fmt"
)

func select0() {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)

	ch1 <- 1
	ch2 <- 1

	select {
	case val := <-ch1:
		fmt.Println("ch1 val", val)
	case ch2 <- 1:
		fmt.Println("put val to ch2")
	default:
		fmt.Println("default case")
	}
	fmt.Println("finish")
}

func select1() {
	ch1 := make(chan int, 2)
	ch1 <- 1
	ch1 <- 2
	ch2 := make(chan int, 2)
	ch2 <- 3

LOOP:
	for {
		select {
		case v1 := <-ch1:
			fmt.Println("chan1 val", v1)
		case v2 := <-ch2:
			fmt.Println("chan2 val", v2)
		default:
			break LOOP
		}
	}
}

func select2() {
	ch := make(chan int, 3)

	ch <- 1
	ch <- 2
	ch <- 3

	close(ch)
	fmt.Println("channel closed")

	// for v := range ch {
	// 	fmt.Println("read:", v)
	// }
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				fmt.Println("select: channel empty and closed — exit")
				return
			}
			fmt.Println("select: read value", v)
		}
	}
}

func select3() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go func() {
		ch2 <- 0
		close(ch2)
		ch1 <- 1
		ch1 <- 2
		ch1 <- 3
		ch1 <- 4
		ch1 <- 5
		ch1 <- 6
		close(ch1)
	}()

	for {
		select {
		case v1, ok := <-ch1:
			if !ok {
				fmt.Println("ch1 closed")
				ch1 = nil
			}
			fmt.Println("chan1 val", v1)
		case v2, ok := <-ch2:
			if !ok {
				fmt.Println("ch2 closed")
				ch2 = nil
			}
			fmt.Println("chan2 val", v2)
		}
		if ch1 == nil && ch2 == nil {
			fmt.Print("Two channels closed")
			break
		}
	}
}

func main() {
	// select0()
	// select1()
	// select2()
	select3()
}

