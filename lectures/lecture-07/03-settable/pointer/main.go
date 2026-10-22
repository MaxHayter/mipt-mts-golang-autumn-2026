// Менять можно только то, до чего есть адрес: копия против указателя.
//
//	go run ./code/03-settable/pointer
package main

import (
	"fmt"
	"reflect"
)

type Config struct {
	Port int
}

// try выполняет f и печатает панику вместо падения программы.
func try(name string, f func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%-24s PANIC: %v\n", name, r)
		}
	}()
	f()
	fmt.Printf("%-24s ok\n", name)
}

func main() {
	cfg := Config{Port: 80}

	// ValueOf получил копию cfg — менять копию бессмысленно, reflect запрещает.
	byValue := reflect.ValueOf(cfg)
	fmt.Println("по значению CanSet:", byValue.Field(0).CanSet())
	try("SetInt через копию", func() { byValue.Field(0).SetInt(8080) })

	// Указатель + Elem: теперь Value смотрит на сам cfg.
	byPtr := reflect.ValueOf(&cfg).Elem()
	fmt.Println("\nчерез указатель CanSet:", byPtr.Field(0).CanSet())
	try("SetInt через указатель", func() { byPtr.Field(0).SetInt(8080) })
	fmt.Println("cfg.Port =", cfg.Port)
}
