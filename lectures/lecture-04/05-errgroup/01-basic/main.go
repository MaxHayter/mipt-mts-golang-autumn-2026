package main

import (
	"fmt"

	"golang.org/x/sync/errgroup"
)

type field struct{ name, value string }

func validate(f field) error {
	if f.value == "" {
		return fmt.Errorf("%s: пустое значение", f.name)
	}
	fmt.Printf("%s: ок\n", f.name)
	return nil
}

// Три независимые проверки не должны ждать друг друга по очереди.
// errgroup.Group — это просто WaitGroup, который попутно запоминает
// первую ошибку; контекст ему для этого не нужен.
func main() {
	var g errgroup.Group

	fields := []field{
		{"имя", "Аня"},
		{"email", "a@example.com"},
		{"телефон", ""},
	}
	for _, f := range fields {
		g.Go(func() error { return validate(f) })
	}

	if err := g.Wait(); err != nil {
		fmt.Println("итог: первая ошибка —", err)
	}
}
