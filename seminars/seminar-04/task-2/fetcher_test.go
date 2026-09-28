package fetcher

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/goleak"
)

type MockServerResponse struct {
	data  []string
	err   error
	delay time.Duration
}

type MockFetcher struct {
	response map[string]MockServerResponse
}

func NewFetcherMock() *MockFetcher {
	return &MockFetcher{
		response: make(map[string]MockServerResponse),
	}
}

func (mf MockFetcher) SetResponse(server string, data []string, err error, delay time.Duration) {
	mf.response[server] = MockServerResponse{
		data:  data,
		err:   err,
		delay: delay,
	}
}

func (mf MockFetcher) Get(server string, query string) ([]string, error) {

	resp, ok := mf.response[server]
	if !ok {
		return nil, errors.New("Unknown server")
	}

	if resp.delay > 0 {
		time.Sleep(resp.delay)
	}

	return resp.data, resp.err
}

func TestFirstSuccessResult(t *testing.T) {
	defer goleak.VerifyNone(t)

	mock := NewFetcherMock()
	mock.SetResponse("server_1", []string{"response_1"}, nil, 10*time.Millisecond)
	mock.SetResponse("server_2", []string{"response_2"}, nil, 20*time.Millisecond)
	mock.SetResponse("server_3", []string{"response_3"}, nil, 30*time.Millisecond)

	ctx := context.Background()
	resp, err := FetchServersPack(ctx, mock, []string{"server_1", "server_2", "server_3"}, "query")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(resp) == 0 {
		t.Error("expected non-empty result")
	}

	if resp[0] != "response_1" {
		t.Errorf("expect: response_1 got: %v", resp[0])
	}
}
