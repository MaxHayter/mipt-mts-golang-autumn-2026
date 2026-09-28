package main

import "fmt"

type User struct {
	ID    int64
	Email string
}

// Escape - указатель на переменную возвраащется из функции
func escapeCase1() *User {
	u := User{ID: 1, Email: "test@example.com"}
	return &u
}

// Нет escape - из функции возвращается копия переменной
func noEscapeCase1() User {
	u := User{ID: 1, Email: "test@example.com"}
	return u
}

var globalUser *User

// Escape - указатель на локальную переменную присваивается глобальной
func escapeCase2() {
	u := User{ID: 2}
	globalUser = &u
}

// Нет escape - в канал отправляется копия локальной переменной
func noEscapeCase2(ch chan User) {
	u := User{ID: 3}
	ch <- u
}

// Escape - передача указателя в неизвестную функцию
func escapeCase3(fn func(*User)) {
	u := User{ID: 4}
	fn(&u)
}

// Нет escape - переменная используется только локально
func noEscapeCase3() {
	u := User{ID: 5}
	processUserLocally(&u)
}

// Нет escape - переменная используется только локально
func processUserLocally(u *User) {
	fmt.Println(u.ID)
}

// Escape - слайс указателей возвращается из функции
func escapeCase4() []*User {
	u := User{ID: 6}
	return []*User{&u}
}

// Нет escape - слайс локален
func noEscapeCase4() {
	u := User{ID: 7}
	users := []*User{&u}
	_ = users // Используем локально
}

// Escape - interface{} размещается на heap
func escapeCase5() interface{} {
	u := User{ID: 8}
	return &u // interface{} требует выделения памяти
}

func main() {
	escapeCase1()
	_ = noEscapeCase1()
	escapeCase2()
	escapeCase3(func(u *User) { fmt.Println(u) })
	noEscapeCase3()
	escapeCase4()
	noEscapeCase4()
	escapeCase5()
}
