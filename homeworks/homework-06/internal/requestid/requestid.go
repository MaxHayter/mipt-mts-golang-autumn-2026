// Package requestid передаёт идентификатор запроса через context.
package requestid

import "context"

type key struct{}

func WithContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, key{}, id)
}

// FromContext используется клиентскими интерсепторами для gRPC metadata.
func FromContext(ctx context.Context) string {
	id, _ := ctx.Value(key{}).(string)
	return id
}
