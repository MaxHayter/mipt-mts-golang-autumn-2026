package main

import (
	"fmt"
)

func main() {
	var count int

	for range 1_000 {
		count++
	}

	fmt.Println("count:", count)
}
