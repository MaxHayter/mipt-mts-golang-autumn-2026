// Грабли: встроенная структура — NumField видит Base, а не ID и Created.
//
//	go run ./code/05-pitfalls/embedded
package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

type Base struct {
	ID      int       `json:"id"`
	Created time.Time `json:"created"`
}

// Order встраивает Base: поля ID и Created продвигаются наверх.
type Order struct {
	Base
	Total int `json:"total"`
}

func main() {
	o := Order{Base: Base{ID: 7, Created: time.Now()}, Total: 990}
	t := reflect.TypeOf(o)
	fmt.Println("Order: NumField =", t.NumField())
	for f := range t.Fields() {
		f := f
		fmt.Printf("  %-7s Anonymous=%v\n", f.Name, f.Anonymous)
	}
	fmt.Println("VisibleFields:")
	for _, f := range reflect.VisibleFields(t) {
		fmt.Printf("  %-7s index=%v\n", f.Name, f.Index)
	}
	// encoding/json раскрывает продвинутые поля — свой обход должен делать так же.
	b, _ := json.Marshal(o)
	fmt.Println("json.Marshal:", string(b))
}
