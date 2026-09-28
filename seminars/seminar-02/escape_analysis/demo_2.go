package main

import (
	"testing"
)

type TestUser struct {
	ID    int64
	Email string
	Age   int32
}

// Escapes на heap - медленнее
func createUserEscape() *TestUser {
	u := TestUser{ID: 1, Email: "test@example.com", Age: 25}
	return &u
}

// Остаётся на stack - быстрее
func createUserNoEscape() TestUser {
	u := TestUser{ID: 1, Email: "test@example.com", Age: 25}
	return u
}

// Escapes в slice - медленнее
func createUsersEscapeSlice() []*TestUser {
	u1 := TestUser{ID: 1}
	u2 := TestUser{ID: 2}
	return []*TestUser{&u1, &u2}
}

// Слайс по значению - быстрее
func createUsersNoEscapeSlice() []TestUser {
	u1 := TestUser{ID: 1}
	u2 := TestUser{ID: 2}
	return []TestUser{u1, u2}
}

func BenchmarkEscape(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = createUserEscape()
	}
}

func BenchmarkNoEscape(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = createUserNoEscape()
	}
}

func BenchmarkEscapeSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = createUsersEscapeSlice()
	}
}

func BenchmarkNoEscapeSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = createUsersNoEscapeSlice()
	}
}
