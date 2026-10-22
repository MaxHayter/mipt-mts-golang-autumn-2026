// Даже через указатель: скрытое поле не записать, тип проверяется в рантайме.
//
//	go run ./code/03-settable/limits
package main

import (
	"fmt"
	"reflect"
)

type Config struct {
	Port int
	host string
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
	v := reflect.ValueOf(&cfg).Elem()

	// Неэкспортируемое поле адресуемо, но не устанавливаемо.
	host := v.FieldByName("host")
	fmt.Printf("host: CanAddr=%v CanSet=%v\n", host.CanAddr(), host.CanSet())
	try("SetString в host", func() { host.SetString("evil") })

	// Тип значения проверяется в рантайме, а не компилятором.
	try("SetString в int-поле", func() { v.Field(0).SetString("8080") })
	try("Set(ValueOf(int64))", func() { v.Field(0).Set(reflect.ValueOf(int64(1))) })
}
