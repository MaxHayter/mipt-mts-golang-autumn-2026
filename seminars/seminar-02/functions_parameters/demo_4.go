package demo

type User struct {
	ID      int64
	Active  bool
	Email   string
	Age     int32
	Premium bool
}

// Оптимизированная версия
type UserOpt struct {
	ID      int64
	Email   string
	Age     int32
	Active  bool
	Premium bool
}

// Пример 1: Передача по значению
//
//go:noinline
func processValue(u User) int {
	u.Active = true // не влияет на исходный объект
	return int(u.ID)
}

// Пример 2: Передача по указателю
//
//go:noinline
func processPointer(u *User) int {
	u.Active = true // влияет на исходный объект
	return int(u.ID)
}

// Пример 3: Большая структура
type LargeData struct {
	Values [1000]int64
}

// Плохо: копирует весь массив
//
//go:noinline
func badFunction(d LargeData) int64 {
	return d.Values[0]
}

// Хорошо: копирует только адрес
//
//go:noinline
func goodFunction(d *LargeData) int64 {
	return d.Values[0]
}

type BigItem struct {
	A, B, C, D, E, F, G, H [16]int64
}

//go:noinline
func rangeByValue(items []BigItem) int64 {
	var sum int64
	for _, v := range items {
		sum += v.A[0]
	}
	return sum
}

//go:noinline
func rangeByPointer(items []BigItem) int64 {
	var sum int64
	for i := 0; i < len(items); i++ {
		sum += items[i].A[0]
	}
	return sum
}
