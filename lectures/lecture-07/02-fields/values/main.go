// Значения полей через итератор Value.Fields (Go 1.26).
//
//	go run ./code/02-fields/values
package main

import (
	"fmt"
	"reflect"
)

type UserID int

type User struct {
	ID       UserID `json:"id"`
	Login    string `json:"login" validate:"required,min=3"`
	RealName string `json:"-"`
	Email    string `json:"email,omitempty"`
	password string // неэкспортируемое
}

func main() {
	u := User{ID: 42, Login: "rob", RealName: "Rob Pike", password: "secret"}

	// Пары (описание поля, значение) — обычный for range.
	for f, fv := range reflect.ValueOf(u).Fields() {
		if !f.IsExported() {
			// fv.Interface() здесь паникует: reflect не выдаёт наружу скрытые поля.
			fmt.Printf("%-8s <скрыто, CanInterface=%v>\n", f.Name, fv.CanInterface())
			continue
		}
		fmt.Printf("%-8s %v\n", f.Name, fv.Interface())
	}
}
