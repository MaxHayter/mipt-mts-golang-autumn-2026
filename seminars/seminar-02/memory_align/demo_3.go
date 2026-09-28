package main

import (
	"fmt"
	"unsafe"
)

// Пример 1: Размеры типов
func demoSizes() {
	var a int64
	var b int32
	var c int16
	var d int8

	fmt.Printf("\nint64: %d bytes\n", unsafe.Sizeof(a))
	fmt.Printf("int32: %d bytes\n", unsafe.Sizeof(b))
	fmt.Printf("int32: %d bytes\n", unsafe.Sizeof(c))
	fmt.Printf("int32: %d bytes\n\n", unsafe.Sizeof(d))
	fmt.Printf("Alignment int64: %d\n\n", unsafe.Alignof(a))
}

// Пример 2: Выравнивание в структуре (bad design)
type BadStruct struct {
	A int8  // 1 byte
	B int64 // 8 bytes
	C int16 // 2 bytes
}

// Пример 3: Оптимизированная структура (good design)
type GoodStruct struct {
	B int64 // 8 bytes
	A int8  // 1 byte
	C int16 // 2 bytes
}

func main() {
	demoSizes()
	fmt.Printf("BadStruct size: %d bytes\n", unsafe.Sizeof(BadStruct{}))
	fmt.Printf("GoodStruct size: %d bytes\n", unsafe.Sizeof(GoodStruct{}))
}
