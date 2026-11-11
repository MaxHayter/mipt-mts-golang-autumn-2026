// Package incident — заготовка адаптера gRPC-клиента. Реализация — часть ДЗ.
package incident

import (
	"context"

	"homework06/internal/contract"
	httpapi "homework06/internal/transport/http"
)

type Client struct{}

// New должен создать одно переиспользуемое gRPC-соединение и сгенерированный клиент.
func New(address string) (*Client, error) {
	// TODO: подключение к address, клиентские interceptors.
	return &Client{}, nil
}

func (c *Client) Close() error {
	// TODO: закрыть gRPC-соединение.
	return nil
}

func notImplemented() error {
	return &contract.Error{Code: "Unavailable", Message: "grpc client is not implemented"}
}

func (c *Client) Create(ctx context.Context, title, service string) (*contract.Incident, error) {
	// TODO: вызов gRPC, преобразование ответа и status.Code(err).
	return nil, notImplemented()
}

func (c *Client) Get(ctx context.Context, id string) (*contract.Incident, error) {
	// TODO: вызов gRPC, преобразование ответа и status.Code(err).
	return nil, notImplemented()
}

func (c *Client) Assign(ctx context.Context, id, user string) (*contract.Incident, error) {
	// TODO: вызов gRPC, преобразование ответа и status.Code(err).
	return nil, notImplemented()
}

func (c *Client) ChangeStatus(ctx context.Context, id, status string) (*contract.Incident, error) {
	// TODO: вызов gRPC, преобразование ответа и status.Code(err).
	return nil, notImplemented()
}

func (c *Client) AddComment(ctx context.Context, id, author, text string) (*contract.Incident, error) {
	// TODO: вызов gRPC, преобразование ответа и status.Code(err).
	return nil, notImplemented()
}

func (c *Client) Watch(ctx context.Context, id string) (httpapi.Stream, error) {
	// TODO: открыть gRPC-поток с переданным ctx и обернуть его Recv.
	// Обрабатывайте ошибки как при открытии потока, так и при Recv.
	return nil, notImplemented()
}

var _ httpapi.Client = (*Client)(nil)
