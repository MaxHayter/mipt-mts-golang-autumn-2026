package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func fetch(ctx context.Context, name string, delay time.Duration, fail bool) error {
	select {
	case <-time.After(delay):
	case <-ctx.Done():
		fmt.Printf("%s: отменено (%v)\n", name, ctx.Err())
		return ctx.Err()
	}
	if fail {
		return fmt.Errorf("%s: %w", name, errors.New("сервис недоступен"))
	}
	fmt.Printf("%s: получено\n", name)
	return nil
}

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error { return fetch(ctx, "профиль", 30*time.Millisecond, false) })
	g.Go(func() error { return fetch(ctx, "баланс", 50*time.Millisecond, true) })
	g.Go(func() error { return fetch(ctx, "рекомендации", 200*time.Millisecond, false) })

	if err := g.Wait(); err != nil {
		fmt.Println("итог: первая ошибка —", err)
	}
}
