package main

import (
	"fmt"
	"runtime"
)

func main() {
	go fmt.Println("Hello, GO!")
	runtime.Gosched()
}
