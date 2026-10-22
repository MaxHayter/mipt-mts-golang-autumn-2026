// Грабли: в any ничего не положили или пришёл nil-указатель.
//
//	go run ./code/05-pitfalls/nilvalue
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
	// Пустой any: Value невалидный, у него нет даже Kind.
	nothing := reflect.ValueOf(nil)
	fmt.Println("ValueOf(nil): IsValid =", nothing.IsValid(), "| Kind =", nothing.Kind())
	try("NumField на пустом", func() { nothing.NumField() })

	// nil-указатель на структуру: Elem() даёт невалидный Value.
	var cfg *Config
	ptr := reflect.ValueOf(cfg)
	fmt.Println("\n(*Config)(nil): Kind =", ptr.Kind(), "| IsNil =", ptr.IsNil())
	try("NumField на пустом", func() { ptr.NumField() })
	elem := ptr.Elem()
	fmt.Println("Elem().IsValid =", elem.IsValid())
	try("Field(0) через nil", func() { elem.Field(0) })
}
