// Самый короткий путь в кодогенерацию: stringer пишет метод String() за нас.
//
//	go generate ./code/07-stringer   // создаёт status_string.go
//	go run ./code/07-stringer
package main

import "fmt"

//go:generate go tool stringer -type=Status -trimprefix=Status

type Status int

const (
	StatusNew Status = iota
	StatusPaid
	// StatusRefunded
	StatusShipped
	StatusCancelled
)

func main() {
	s := StatusShipped
	fmt.Println(s) // Shipped — fmt вызывает String()
	fmt.Printf("%d %v\n", s, s)
	fmt.Println(Status(42)) // Status(42) — неизвестное значение не молчит

	// Вставьте StatusRefunded между StatusPaid и StatusShipped и запустите
	// go run без go generate: сборка упадёт — сгенерированный файл проверяет,
	// что значения констант не сдвинулись.
}
