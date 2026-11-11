package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"homework06/internal/contract"
	"homework06/internal/requestid"
)

type fakeClient struct {
	call  func(context.Context, string, ...string) (*contract.Incident, error)
	watch func(context.Context, string) (Stream, error)
}

func (f fakeClient) Create(c context.Context, a, b string) (*contract.Incident, error) {
	return f.call(c, "Create", a, b)
}
func (f fakeClient) Get(c context.Context, a string) (*contract.Incident, error) {
	return f.call(c, "Get", a)
}
func (f fakeClient) Assign(c context.Context, a, b string) (*contract.Incident, error) {
	return f.call(c, "Assign", a, b)
}
func (f fakeClient) ChangeStatus(c context.Context, a, b string) (*contract.Incident, error) {
	return f.call(c, "ChangeStatus", a, b)
}
func (f fakeClient) AddComment(c context.Context, a, b, d string) (*contract.Incident, error) {
	return f.call(c, "AddComment", a, b, d)
}
func (f fakeClient) Watch(c context.Context, id string) (Stream, error) { return f.watch(c, id) }

type fakeStream struct {
	recv func() (*contract.Incident, error)
}

func (f fakeStream) Recv() (*contract.Incident, error) { return f.recv() }

func TestRoutes(t *testing.T) {
	cases := []struct {
		name, method, path, body, action string
		args                             []string
		code                             int
	}{
		{"create", "POST", "/incidents", `{"title":"Сбой","service":"payments"}`, "Create", []string{"Сбой", "payments"}, 201},
		{"get", "GET", "/incidents/abc", "", "Get", []string{"abc"}, 200},
		{"assign", "PUT", "/incidents/abc/assignee", `{"user":"Nick"}`, "Assign", []string{"abc", "Nick"}, 200},
		{"status", "PUT", "/incidents/abc/status", `{"status":"OPEN"}`, "ChangeStatus", []string{"abc", "OPEN"}, 200},
		{"comment", "POST", "/incidents/abc/comments", `{"author":"Nick","text":"test"}`, "AddComment", []string{"abc", "Nick", "test"}, 200},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			client := fakeClient{call: func(ctx context.Context, action string, args ...string) (*contract.Incident, error) {
				called = true
				if action != tc.action || strings.Join(args, "|") != strings.Join(tc.args, "|") {
					t.Fatalf("action=%s args=%v", action, args)
				}
				if requestid.FromContext(ctx) != "test-id" {
					t.Error("request ID lost")
				}
				if _, ok := ctx.Deadline(); !ok {
					t.Error("deadline missing")
				}
				return &contract.Incident{ID: "abc", Version: 1, Status: "OPEN"}, nil
			}}
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			req.Header.Set("X-Request-ID", "test-id")
			w := httptest.NewRecorder()
			New(client, time.Second).ServeHTTP(w, req)
			if !called || w.Code != tc.code {
				t.Fatalf("called=%v code=%d body=%s", called, w.Code, w.Body)
			}
			if !strings.Contains(w.Body.String(), `"comments":[]`) || w.Header().Get("X-Request-ID") != "test-id" {
				t.Fatalf("response=%s headers=%v", w.Body, w.Header())
			}
		})
	}
}

func TestErrors(t *testing.T) {
	cases := []struct {
		name    string
		err     error
		code    int
		message string
	}{
		{"invalid", &contract.Error{Code: "InvalidArgument", Message: "bad input"}, 400, "bad input"},
		{"missing", &contract.Error{Code: "NotFound", Message: "not found"}, 404, "not found"},
		{"conflict", &contract.Error{Code: "FailedPrecondition", Message: "closed"}, 409, "closed"},
		{"unavailable", &contract.Error{Code: "Unavailable", Message: "unavailable"}, 503, "unavailable"},
		{"deadline", context.DeadlineExceeded, 504, "request timeout"},
		{"internal", errors.New("secret database details"), 500, "internal server error"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := fakeClient{call: func(context.Context, string, ...string) (*contract.Incident, error) { return nil, tc.err }}
			w := httptest.NewRecorder()
			New(c, time.Second).ServeHTTP(w, httptest.NewRequest("GET", "/incidents/x", nil))
			var body map[string]string
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if w.Code != tc.code || body["error"] != tc.message || body["request_id"] == "" {
				t.Fatalf("code=%d body=%v", w.Code, body)
			}
		})
	}
}

func TestMalformedJSON(t *testing.T) {
	for _, body := range []string{`{`, `{} {}`, `[]`, `null`, `"text"`, `{"title":12}`} {
		t.Run(body, func(t *testing.T) {
			c := fakeClient{call: func(context.Context, string, ...string) (*contract.Incident, error) {
				t.Fatal("backend called for malformed JSON")
				return nil, nil
			}}
			w := httptest.NewRecorder()
			New(c, time.Second).ServeHTTP(w, httptest.NewRequest("POST", "/incidents", strings.NewReader(body)))
			if w.Code != 400 {
				t.Fatalf("code=%d", w.Code)
			}
		})
	}
}

func TestSSEFirstError(t *testing.T) {
	c := fakeClient{watch: func(context.Context, string) (Stream, error) {
		return fakeStream{recv: func() (*contract.Incident, error) { return nil, &contract.Error{Code: "NotFound", Message: "missing"} }}, nil
	}}
	w := httptest.NewRecorder()
	New(c, time.Second).ServeHTTP(w, httptest.NewRequest("GET", "/incidents/x/events", nil))
	if w.Code != 404 || !strings.HasPrefix(w.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("code=%d headers=%v", w.Code, w.Header())
	}
}

func TestSSEFinalSnapshot(t *testing.T) {
	calls := 0
	var streamCtx context.Context
	c := fakeClient{watch: func(ctx context.Context, id string) (Stream, error) {
		streamCtx = ctx
		if _, ok := ctx.Deadline(); ok {
			t.Error("unary deadline applied to stream")
		}
		return fakeStream{recv: func() (*contract.Incident, error) {
			calls++
			if calls == 1 {
				return &contract.Incident{ID: id, Version: 4, Status: "RESOLVED"}, nil
			}
			return nil, io.EOF
		}}, nil
	}}
	w := httptest.NewRecorder()
	New(c, time.Second).ServeHTTP(w, httptest.NewRequest("GET", "/incidents/x/events", nil))
	if !w.Flushed || !strings.Contains(w.Body.String(), "event: snapshot\n") || !strings.Contains(w.Body.String(), `"status":"RESOLVED"`) {
		t.Fatalf("body=%s", w.Body)
	}
	if streamCtx.Err() == nil {
		t.Fatal("subscription context not canceled on return")
	}
}

func TestSSEDisconnectCancelsBackend(t *testing.T) {
	canceled := make(chan struct{})
	c := fakeClient{watch: func(ctx context.Context, id string) (Stream, error) {
		first := true
		return fakeStream{recv: func() (*contract.Incident, error) {
			if first {
				first = false
				return &contract.Incident{ID: id, Version: 1, Status: "OPEN"}, nil
			}
			<-ctx.Done()
			close(canceled)
			return nil, ctx.Err()
		}}, nil
	}}
	srv := httptest.NewServer(New(c, time.Second))
	defer srv.Close()
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(srv.URL + "/incidents/x/events")
	if err != nil {
		t.Fatal(err)
	}
	if err := resp.Body.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case <-canceled:
	case <-time.After(3 * time.Second):
		t.Fatal("HTTP disconnect did not cancel backend stream")
	}
}

func TestSerializationErrorKeepsRequestID(t *testing.T) {
	c := fakeClient{call: func(context.Context, string, ...string) (*contract.Incident, error) {
		return &contract.Incident{Comments: []contract.Comment{{CreatedAt: time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC)}}}, nil
	}}
	req := httptest.NewRequest("GET", "/incidents/x", nil)
	req.Header.Set("X-Request-ID", "serialize-test")
	w := httptest.NewRecorder()
	New(c, time.Second).ServeHTTP(w, req)
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if w.Code != 500 || body["request_id"] != "serialize-test" || body["error"] != "internal server error" {
		t.Fatalf("code=%d body=%v", w.Code, body)
	}
}

func TestSSEErrorAfterSnapshot(t *testing.T) {
	calls := 0
	c := fakeClient{watch: func(context.Context, string) (Stream, error) {
		return fakeStream{recv: func() (*contract.Incident, error) {
			calls++
			if calls == 1 {
				return &contract.Incident{ID: "x", Version: 1, Status: "OPEN"}, nil
			}
			return nil, &contract.Error{Code: "Internal", Message: "secret details"}
		}}, nil
	}}
	w := httptest.NewRecorder()
	New(c, time.Second).ServeHTTP(w, httptest.NewRequest("GET", "/incidents/x/events", nil))
	body := w.Body.String()
	if w.Code != 200 || !strings.Contains(body, "event: snapshot\n") || !strings.Contains(body, "event: error\n") || !strings.Contains(body, `"request_id":`) || strings.Contains(body, "secret details") {
		t.Fatalf("code=%d body=%s", w.Code, body)
	}
}

func TestSSEInvalidBackendResponse(t *testing.T) {
	cases := []struct {
		name   string
		stream Stream
	}{
		{"nil_stream", nil},
		{"nil_snapshot", fakeStream{recv: func() (*contract.Incident, error) { return nil, nil }}},
		{"early_eof", fakeStream{recv: func() (*contract.Incident, error) { return nil, io.EOF }}},
		{"invalid_time", fakeStream{recv: func() (*contract.Incident, error) {
			return &contract.Incident{Comments: []contract.Comment{{CreatedAt: time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC)}}}, nil
		}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := fakeClient{watch: func(context.Context, string) (Stream, error) { return tc.stream, nil }}
			w := httptest.NewRecorder()
			New(c, time.Second).ServeHTTP(w, httptest.NewRequest("GET", "/incidents/x/events", nil))
			if w.Code != 500 || !strings.HasPrefix(w.Header().Get("Content-Type"), "application/json") {
				t.Fatalf("code=%d body=%s", w.Code, w.Body)
			}
		})
	}
}

func TestUnaryDeadlineAndCancellation(t *testing.T) {
	c := fakeClient{call: func(ctx context.Context, _ string, _ ...string) (*contract.Incident, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	}}
	w := httptest.NewRecorder()
	New(c, time.Millisecond).ServeHTTP(w, httptest.NewRequest("GET", "/incidents/x", nil))
	if w.Code != 504 {
		t.Fatalf("code=%d body=%s", w.Code, w.Body)
	}
}
