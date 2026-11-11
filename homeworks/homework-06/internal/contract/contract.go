// Package contract содержит обычные Go-типы на границе HTTP и gRPC-клиента.
// Это не protobuf-контракт и не реализация бизнес-логики.
package contract

import "time"

type Comment struct {
	Author    string    `json:"author"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type Incident struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Service  string    `json:"service"`
	Status   string    `json:"status"`
	Assignee string    `json:"assignee"`
	Comments []Comment `json:"comments"`
	Version  uint64    `json:"version"`
}

// Error возвращает адаптер gRPC-клиента, преобразуя status.Code(err).
// Code: InvalidArgument, NotFound, FailedPrecondition, Unavailable,
// DeadlineExceeded или Internal. Message — безопасный текст для пользователя.
type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string { return e.Message }
