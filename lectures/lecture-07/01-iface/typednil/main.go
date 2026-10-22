// Разгадка typed nil: в интерфейсе заполнена половина «тип».
//
//	go run ./code/01-iface/typednil
package main

import (
	"errors"
	"fmt"
	"reflect"
)

type UserID int

type NotFoundError struct{ ID UserID }

func (e *NotFoundError) Error() string { return fmt.Sprintf("user %d not found", e.ID) }

// find — та самая ловушка из лекции 2: err типа *NotFoundError,
// при успехе он nil, но возвращается как error — а это уже пара (тип, nil).
func find(id UserID) error {
	var err *NotFoundError
	if id == 0 {
		err = &NotFoundError{ID: id}
	}
	return err
}

func main() {
	err := find(1)
	fmt.Println("err == nil:", err == nil)
	fmt.Println("тип внутри:", reflect.TypeOf(err))
	fmt.Println("значение внутри nil:", reflect.ValueOf(err).IsNil())

	var nf *NotFoundError
	fmt.Println("errors.As всё равно находит тип:", errors.As(err, &nf))

	// Настоящий пустой интерфейс: обе половины пусты.
	var empty error
	fmt.Println("\nTypeOf(nil error):", reflect.TypeOf(empty))
	fmt.Println("ValueOf(nil error).IsValid():", reflect.ValueOf(empty).IsValid())
}
