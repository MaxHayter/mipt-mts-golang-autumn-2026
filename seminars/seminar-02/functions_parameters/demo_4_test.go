package demo

import (
	"testing"
)

var Sink int64

// Бенчмарк: передача по значению
func BenchmarkByValue(b *testing.B) {
	u := User{ID: 1, Email: "test@example.com"}
	var s int64
	for i := 0; i < b.N; i++ {
		s += int64(processValue(u))
	}
	Sink = s
}

// Бенчмарк: передача по указателю
func BenchmarkByPointer(b *testing.B) {
	u := User{ID: 1, Email: "test@example.com"}
	var s int64
	for i := 0; i < b.N; i++ {
		s += int64(processPointer(&u))
	}
	Sink = s
}

// Бенчмарк: большая структура по значению
func BenchmarkLargeByValue(b *testing.B) {
	d := LargeData{}
	var s int64
	for i := 0; i < b.N; i++ {
		s += badFunction(d)
	}
	Sink = s
}

// Бенчмарк: большая структура по указателю
func BenchmarkLargeByPointer(b *testing.B) {
	d := LargeData{}
	var s int64
	for i := 0; i < b.N; i++ {
		s += goodFunction(&d)
	}
	Sink = s
}

// Бенчмарк: range по слайсу структур (копирование каждого элемента)
func BenchmarkRangeStructByValue(b *testing.B) {
	items := make([]BigItem, 1000)
	var sum int64
	for i := 0; i < b.N; i++ {
		sum += rangeByValue(items)
	}
	Sink = sum
}

// Бенчмарк: проход по слайсу структур по индексу
func BenchmarkRangeStructByIndex(b *testing.B) {
	items := make([]BigItem, 1000)
	var sum int64
	for i := 0; i < b.N; i++ {
		sum += rangeByPointer(items)
	}
	Sink = sum
}
