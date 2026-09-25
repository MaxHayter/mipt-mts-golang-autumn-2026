package main

import (
	"context"
	"fmt"
)

// Свой неэкспортируемый тип ключа: даже если другой пакет заведёт
// собственный requestIDKey с тем же именем, столкновения не будет —
// context.Value сравнивает интерфейсы целиком, тип и значение сразу,
// а этот тип принадлежит только этому пакету.
type ctxKey int

const requestIDKey ctxKey = 0

// Ключом необязательно быть числом или строкой — годится любой
// сравнимый тип, в том числе пустая структура. Сравнивать у нее
// нечего, кроме самого типа, поэтому она надёжна как ключ сама по
// себе, без отдельной константы вроде requestIDKey.
type sessionKey struct{}

func handle(ctx context.Context) {
	logStep(ctx, "начали обработку")
	logStep(ctx, "закончили обработку")
}

// requestID и session не передаются отдельными параметрами через всю
// цепочку вызовов — оба приехали вместе с ctx.
func logStep(ctx context.Context, msg string) {
	id, _ := ctx.Value(requestIDKey).(string)
	session, _ := ctx.Value(sessionKey{}).(string)
	fmt.Printf("[%s/%s] %s\n", id, session, msg)
}

func main() {
	ctx := context.WithValue(context.Background(), requestIDKey, "req-42")
	ctx = context.WithValue(ctx, sessionKey{}, "sess-7")
	handle(ctx)

	_, ok := context.Background().Value(requestIDKey).(string)
	fmt.Println("значение есть в пустом контексте:", ok)
}
