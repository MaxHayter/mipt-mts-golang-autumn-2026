// Вход для main.go: он читает этот файл как текст и разбирает парсером.
package main

import "time"

type UserID int64

//gen:json
type User struct {
	ID      UserID    `json:"id"`
	Login   string    `json:"login"`
	Tags    []string  `json:"tags"`
	Created time.Time `json:"created"`
}
