// Обход полей структуры и чтение тегов — то, что делает encoding/json внутри.
//
//	go run ./code/02-fields/tags
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
	t := reflect.TypeFor[User]()
	fmt.Printf("%v: %d полей\n", t, t.NumField())

	// Классическая форма — так написан почти весь существующий код.
	for f := range t.Fields() {
		// f := f
		fmt.Printf("  %-8s type=%-11v kind=%-6v json=%-18q validate=%q exported=%v\n",
			f.Name, f.Type, f.Type.Kind(), f.Tag.Get("json"), f.Tag.Get("validate"), f.IsExported())
	}

	// Тег — просто строка. Формат key:"value" — соглашение, которое понимает Tag.Get.
	f, _ := t.FieldByName("Email")
	name, ok := f.Tag.Lookup("json")
	fmt.Printf("\nсырой тег Email: %s\nLookup(json): %q %v\n", f.Tag, name, ok)
}
