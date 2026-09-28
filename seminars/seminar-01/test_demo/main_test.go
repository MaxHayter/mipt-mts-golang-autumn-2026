package main

import (
	"fmt"
	"testing"
)

// === Простой тест ===
func TestAdd(t *testing.T) {
	result := Add(2, 3)
	expected := 5

	if result != expected {
		t.Errorf("Add(2, 3) = %d, хотели %d", result, expected)
	}
}

// === Table-driven тесты (рекомендуемый подход) ===
func TestAddTableDriven(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 2, 3, 5},
		{"negative numbers", -2, -3, -5},
		{"mixed", -2, 3, 1},
		{"zero", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Add(%d, %d) = %d, хотели %d",
					tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// === Тест с ошибками ===
func TestDivide(t *testing.T) {
	tests := []struct {
		name      string
		a, b      int
		expected  int
		shouldErr bool
	}{
		{"normal division", 10, 2, 5, false},
		{"zero divisor", 10, 0, 0, true},
		{"negative", -10, 2, -5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Divide(tt.a, tt.b)

			if tt.shouldErr {
				if err == nil {
					t.Fatalf("Divide(%d, %d) ожидала ошибку", tt.a, tt.b)
				}
				return
			}

			if err != nil {
				t.Fatalf("Divide(%d, %d) вернула ошибку: %v",
					tt.a, tt.b, err)
			}

			if result != tt.expected {
				t.Errorf("Divide(%d, %d) = %d, хотели %d",
					tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// === Бенчмарк ===
func BenchmarkFibonacci(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fibonacci(20)
	}
}

// === Бенчмарк с таблицей ===
func BenchmarkFibonacciValues(b *testing.B) {
	values := []int{10, 15, 20, 25}

	for _, val := range values {
		b.Run(fmt.Sprintf("n=%d", val), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Fibonacci(val)
			}
		})
	}
}

// === Субтесты с Helper ===
func TestMultipleSubtests(t *testing.T) {
	t.Run("group1", func(t *testing.T) {
		if Add(1, 1) != 2 {
			t.Error("1 + 1 должно быть 2")
		}
	})

	t.Run("group2", func(t *testing.T) {
		result, _ := Divide(4, 2)
		if result != 2 {
			t.Error("4 / 2 должно быть 2")
		}
	})
}
