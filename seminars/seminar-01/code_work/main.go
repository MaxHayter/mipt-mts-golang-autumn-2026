package main

import (
	"fmt"
	"log"
	"os"
)

// CalculateSum складывает два числа
// BUG: эта функция никогда не вызывается
func CalculateSum(a, b int) int {
	return a + b
}

// ProcessFile читает файл и выводит количество строк
func ProcessFile(filename string) (int, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	// Ошибка: data не используется корректно
	_ = data
	return 0, nil
}

func main() {
	// Неправильный формат для fmt.Printf
	x := 42
	fmt.Printf("Значение: %v\n", x, x) // Лишний аргумент

	unused := "это не используется"
	_ = unused

	// Ошибка: переменная не инициализирована перед использованием
	var result int
	log.Println(result)
}
