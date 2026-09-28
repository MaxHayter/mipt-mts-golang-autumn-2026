package main

import "fmt"

// Add складывает два числа
func Add(a, b int) int {
	return a + b
}

// Divide делит a на b
func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("деление на ноль")
	}
	return a / b, nil
}

// Fibonacci возвращает n-ое число Фибоначчи
func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}
