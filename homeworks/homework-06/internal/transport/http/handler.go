// Package httpapi содержит готовый HTTP-слой задания.
package httpapi

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"homework06/internal/contract"
	"homework06/internal/requestid"
)

// Stream реализуется обёрткой над сгенерированным gRPC-потоком.
// Штатное завершение Recv возвращает io.EOF. Отмена ctx в Watch должна
// отменять реальный gRPC-вызов, в том числе заблокированный Recv.
type Stream interface {
	Recv() (*contract.Incident, error)
}

// Client реализуют студенты в internal/client/incident.
// Все методы обращаются к incident-service по gRPC.
type Client interface {
	Create(context.Context, string, string) (*contract.Incident, error)
	Get(context.Context, string) (*contract.Incident, error)
	Assign(context.Context, string, string) (*contract.Incident, error)
	ChangeStatus(context.Context, string, string) (*contract.Incident, error)
	AddComment(context.Context, string, string, string) (*contract.Incident, error)
	Watch(context.Context, string) (Stream, error)
}

type handler struct {
	client  Client
	timeout time.Duration
}

func New(client Client, unaryTimeout time.Duration) http.Handler {
	h := &handler{client: client, timeout: unaryTimeout}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /incidents", h.create)
	mux.HandleFunc("GET /incidents/{id}", h.get)
	mux.HandleFunc("PUT /incidents/{id}/assignee", h.assign)
	mux.HandleFunc("PUT /incidents/{id}/status", h.status)
	mux.HandleFunc("POST /incidents/{id}/comments", h.comment)
	mux.HandleFunc("GET /incidents/{id}/events", h.events)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = rand.Text()
		}
		w.Header().Set("X-Request-ID", id)
		start := time.Now()
		mux.ServeHTTP(w, r.WithContext(requestid.WithContext(r.Context(), id)))
		slog.Info("http request", "request_id", id, "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))
	})
}

func readJSON(w http.ResponseWriter, r *http.Request, out any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	var raw json.RawMessage
	if err := dec.Decode(&raw); err != nil {
		return &contract.Error{Code: "InvalidArgument", Message: "invalid json"}
	}
	if len(bytes.TrimSpace(raw)) == 0 || bytes.TrimSpace(raw)[0] != '{' {
		return &contract.Error{Code: "InvalidArgument", Message: "expected json object"}
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return &contract.Error{Code: "InvalidArgument", Message: "invalid json fields"}
	}
	if err := dec.Decode(new(any)); !errors.Is(err, io.EOF) {
		return &contract.Error{Code: "InvalidArgument", Message: "expected one json object"}
	}
	return nil
}

func publicError(err error) (int, string) {
	var e *contract.Error
	if errors.As(err, &e) {
		switch e.Code {
		case "InvalidArgument":
			return 400, e.Message
		case "NotFound":
			return 404, e.Message
		case "FailedPrecondition":
			return 409, e.Message
		case "Unavailable":
			return 503, e.Message
		case "DeadlineExceeded":
			return 504, e.Message
		}
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return 504, "request timeout"
	}
	return 500, "internal server error"
}

func writeJSON(w http.ResponseWriter, r *http.Request, code int, value any) {
	data, err := json.Marshal(value)
	if err != nil {
		slog.Error("encode response", "request_id", requestid.FromContext(r.Context()), "error", err)
		code = 500
		data, _ = json.Marshal(map[string]string{
			"error": "internal server error", "request_id": requestid.FromContext(r.Context()),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(append(data, '\n')); err != nil {
		slog.Debug("write response", "error", err)
	}
}
func fail(w http.ResponseWriter, r *http.Request, err error) {
	code, message := publicError(err)
	if code == 500 {
		slog.Error("backend request failed", "request_id", requestid.FromContext(r.Context()), "error", err)
	}
	writeJSON(w, r, code, map[string]string{"error": message, "request_id": requestid.FromContext(r.Context())})
}
func normalize(s *contract.Incident) *contract.Incident {
	if s == nil {
		return nil
	}
	copy := *s
	if copy.Comments == nil {
		copy.Comments = []contract.Comment{}
	}
	return &copy
}
func (h *handler) unary(w http.ResponseWriter, r *http.Request, code int, call func(context.Context) (*contract.Incident, error)) {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()
	s, err := call(ctx)
	if err != nil {
		fail(w, r, err)
		return
	}
	if s == nil {
		fail(w, r, errors.New("nil backend response"))
		return
	}
	writeJSON(w, r, code, normalize(s))
}
func (h *handler) create(w http.ResponseWriter, r *http.Request) {
	var b struct {
		Title   string `json:"title"`
		Service string `json:"service"`
	}
	if err := readJSON(w, r, &b); err != nil {
		fail(w, r, err)
		return
	}
	h.unary(w, r, 201, func(ctx context.Context) (*contract.Incident, error) { return h.client.Create(ctx, b.Title, b.Service) })
}
func (h *handler) get(w http.ResponseWriter, r *http.Request) {
	h.unary(w, r, 200, func(ctx context.Context) (*contract.Incident, error) { return h.client.Get(ctx, r.PathValue("id")) })
}
func (h *handler) assign(w http.ResponseWriter, r *http.Request) {
	var b struct {
		User string `json:"user"`
	}
	if err := readJSON(w, r, &b); err != nil {
		fail(w, r, err)
		return
	}
	h.unary(w, r, 200, func(ctx context.Context) (*contract.Incident, error) {
		return h.client.Assign(ctx, r.PathValue("id"), b.User)
	})
}
func (h *handler) status(w http.ResponseWriter, r *http.Request) {
	var b struct {
		Status string `json:"status"`
	}
	if err := readJSON(w, r, &b); err != nil {
		fail(w, r, err)
		return
	}
	h.unary(w, r, 200, func(ctx context.Context) (*contract.Incident, error) {
		return h.client.ChangeStatus(ctx, r.PathValue("id"), b.Status)
	})
}
func (h *handler) comment(w http.ResponseWriter, r *http.Request) {
	var b struct {
		Author string `json:"author"`
		Text   string `json:"text"`
	}
	if err := readJSON(w, r, &b); err != nil {
		fail(w, r, err)
		return
	}
	h.unary(w, r, 200, func(ctx context.Context) (*contract.Incident, error) {
		return h.client.AddComment(ctx, r.PathValue("id"), b.Author, b.Text)
	})
}
