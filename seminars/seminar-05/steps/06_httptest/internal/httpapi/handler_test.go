package httpapi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"http-seminar/internal/tasks"
)

func TestCreateTask(t *testing.T) {
	handler := New(tasks.New(nil), nil, "")
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBufferString(`{"text":"Написать тест"}`))
	recorder := httptest.NewRecorder()

	handler.Routes().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusCreated)
	}
	if got := recorder.Body.String(); got != "{\"id\":1,\"text\":\"Написать тест\"}\n" {
		t.Fatalf("body = %q", got)
	}
}
