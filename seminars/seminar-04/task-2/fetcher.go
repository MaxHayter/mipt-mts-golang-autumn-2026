package fetcher

import (
	"context"
	"errors"
	"sync"
)

// У нас есть метод fetch.Get(server string, query string) ([]string, error)

// Написать метод для опроса списка серверов одним запросом по следующему алгоритму:
// серверы опрашиваются одновременно,
// возвращается первый успешный ответ, при этом остальные ответы игнорируются, в том числе и ошибки,
// если все сервера ответили ошибкой, то наш метод должен возвращать ошибку

type Fetcher interface {
	Get(server string, query string) ([]string, error)
}

func FetchServersPack(ctx context.Context, f Fetcher, servers []string, query string) ([]string, error) {
	if query == "" {
		return nil, nil
	}

	wg := sync.WaitGroup{}
	respCh := make(chan []string, len(servers))

	for _, s := range servers {
		wg.Add(1)
		go func(server string) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
			}

			resp, err := f.Get(server, query)
			// ERROR: !=
			if err == nil {
				respCh <- resp
			}
		}(s)
	}

	go func() {
		wg.Wait()
		close(respCh)
	}()

	select {
	case r, ok := <-respCh:
		if !ok {
			return nil, errors.New("all servers return error")
		}
		return r, nil
	case <-ctx.Done():
		return nil, errors.New("method FetchServersPack was canceld")
	}

}
