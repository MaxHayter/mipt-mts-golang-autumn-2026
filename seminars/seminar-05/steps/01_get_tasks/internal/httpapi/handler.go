// Package httpapi is the first HTTP version: read tasks as JSON.
package httpapi

import (
	"encoding/json"
	"net/http"

	"http-seminar/internal/tasks"
)

type Handler struct{ tasks *tasks.Service }

func New(taskService *tasks.Service) *Handler { return &Handler{tasks: taskService} }

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/tasks", h.tasksJSON)
	return mux
}

func (h *Handler) tasksJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(h.tasks.List())
}
