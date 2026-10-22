// Интерфейсное значение — пара (тип, данные). reflect просто читает эту пару.
//
//	go run ./code/01-iface/kinds
package main

import (
	"fmt"
	"reflect"
)

type UserID int

type NotFoundError struct{ ID UserID }

func (e NotFoundError) Error() string { return fmt.Sprintf("user %d not found", e.ID) }

func main() {
	// Одна и та же переменная any — разные пары (тип, значение).
	for _, v := range []any{42, UserID(42), "42", []int{4, 2}, &NotFoundError{7}} {
		t := reflect.TypeOf(v)
		fmt.Printf("%-18s type=%-22v kind=%v\n", fmt.Sprint(v), t, t.Kind())
	}

	// Type и Kind — разные вещи: имя типа против «формы» в памяти.
	t := reflect.TypeOf(UserID(1))
	fmt.Println("\nUserID:", "Name =", t.Name(), "| PkgPath =", t.PkgPath(), "| Kind =", t.Kind())
}
