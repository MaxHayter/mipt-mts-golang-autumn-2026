package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"time"
)

// ExpensiveCalculation - функция, которая требует вычислений
func ExpensiveCalculation(n int) int {
	sum := 0
	for i := 0; i < n*1000000; i++ {
		sum += i % 7
	}
	return sum
}

// MemoryAllocator - функция с выделением памяти
func MemoryAllocator(size int) [][]int {
	data := make([][]int, 0)
	for i := 0; i < size; i++ {
		row := make([]int, 1000)
		for j := 0; j < 1000; j++ {
			row[j] = rand.Intn(10000)
		}
		data = append(data, row)
	}
	return data
}

// ExpensiveWithGoroutines - функция с горутинами
func ExpensiveWithGoroutines(n int) {
	done := make(chan struct{})

	for i := 0; i < n; i++ {
		go func(id int) {
			time.Sleep(100 * time.Millisecond)
			ExpensiveCalculation(id)
			done <- struct{}{}
		}(i)
	}

	for i := 0; i < n; i++ {
		<-done
	}
}

// RecursiveFunction - рекурсивная функция для демо CPU
func RecursiveFunction(n int) int {
	if n <= 1 {
		return 1
	}
	return RecursiveFunction(n-1) + RecursiveFunction(n-2)
}

func main() {
	// === CPU Profile ===
	cpuFile, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Println("Ошибка при создании cpu.prof:", err)
		return
	}
	defer cpuFile.Close()

	err = pprof.StartCPUProfile(cpuFile)
	if err != nil {
		fmt.Println("Ошибка при старте CPU профиля:", err)
		return
	}
	defer pprof.StopCPUProfile()

	fmt.Println("=== Профилирование CPU ===")
	fmt.Println("Запуск дорогостоящих вычислений...")

	// Нагрузка на CPU
	for i := 0; i < 10; i++ {
		result := ExpensiveCalculation(i + 1)
		fmt.Printf("ExpensiveCalculation(%d) = %d\n", i+1, result)
	}

	// Рекурсия для демо
	for i := 20; i < 28; i++ {
		result := RecursiveFunction(i)
		fmt.Printf("Fibonacci(%d) = %d\n", i, result)
	}

	// Горутины
	fmt.Println("\nЗапуск операций в горутинах...")
	ExpensiveWithGoroutines(5)

	// === Memory Profile ===
	fmt.Println("\n=== Профилирование памяти ===")
	fmt.Println("Выделение памяти...")

	// Выделение памяти
	memData := MemoryAllocator(100)
	fmt.Printf("Выделено массивов: %d\n", len(memData))

	// Принудительная сборка мусора перед снятием профиля
	runtime.GC()
	time.Sleep(100 * time.Millisecond)

	// Снятие memory профиля
	memFile, err := os.Create("mem.prof")
	if err != nil {
		fmt.Println("Ошибка при создании mem.prof:", err)
		return
	}
	defer memFile.Close()

	err = pprof.WriteHeapProfile(memFile)
	if err != nil {
		fmt.Println("Ошибка при записи heap профиля:", err)
		return
	}

	// === Trace ===
	fmt.Println("\n=== Трассирование выполнения ===")
	traceFile, err := os.Create("trace.out")
	if err != nil {
		fmt.Println("Ошибка при создании trace.out:", err)
		return
	}
	defer traceFile.Close()

	err = trace.Start(traceFile)
	if err != nil {
		fmt.Println("Ошибка при старте трассировки:", err)
		return
	}

	fmt.Println("Запись трассировки...")
	ExpensiveWithGoroutines(3)

	trace.Stop()

	// === Итоги ===
	fmt.Println("\n=== Профилирование завершено ===")
	fmt.Println("Профили созданы:")
	fmt.Println("  cpu.prof   - просмотр: go tool pprof -http=:6060 cpu.prof")
	fmt.Println("  mem.prof   - просмотр: go tool pprof -http=:6061 mem.prof")
	fmt.Println("  trace.out  - просмотр: go tool trace trace.out")

	// Информация о памяти
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("\nРазмер кучи (в МБ): %.2f\n", float64(m.Alloc)/1024/1024)
	fmt.Printf("Всего выделено (в МБ): %.2f\n", float64(m.TotalAlloc)/1024/1024)
	fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
}
