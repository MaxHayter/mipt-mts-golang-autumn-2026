// Грабли: поле типа any — Kind поля Interface, до значения ещё один шаг через Elem.
//
//	go run ./code/05-pitfalls/anyfield
package main

import (
	"fmt"
	"reflect"
)

type Event struct {
	Name    string `json:"name"`
	Payload any    `json:"payload"`
}

func main() {
	ev := reflect.ValueOf(Event{Name: "paid", Payload: 42})
	payload := ev.FieldByName("Payload")
	fmt.Println("Payload: Kind =", payload.Kind(), "| Elem().Kind =", payload.Elem().Kind())
}
