// Package bench сравнивает шесть способов получить JSON из одной и той же структуры.
//
//	go generate ./code/09-bench   // создаёт easy_user_easyjson.go
//	go test -bench . -benchmem ./code/09-bench
package bench

//go:generate go tool easyjson -no_std_marshalers easy_user.go

// EasyUser — копия model.User для easyjson: у model.User уже есть
// MarshalJSON от нашего генератора, два метода с одним именем не уживутся.
//
//easyjson:json
type EasyUser struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Name     string `json:"name"`
	Email    string `json:"email,omitempty"`
	Age      int    `json:"age"`
	Admin    bool   `json:"admin"`
	Password string `json:"-"`
}
