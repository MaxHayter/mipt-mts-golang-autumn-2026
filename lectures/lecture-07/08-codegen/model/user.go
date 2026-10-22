// Package model — структура, для которой gen пишет MarshalJSON.
//
//	go generate ./code/08-codegen/model   // создаёт user_json.go
package model

//go:generate go run ../gen $GOFILE

//gen:json
type User struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Name     string `json:"name"`
	Email    string `json:"email,omitempty"`
	Age      int    `json:"age"`
	Admin    bool   `json:"admin"`
	Password string `json:"-"`
}
