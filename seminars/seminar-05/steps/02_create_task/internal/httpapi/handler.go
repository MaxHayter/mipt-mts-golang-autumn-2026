// Package httpapi adds task creation to the first JSON endpoint.
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
	mux.HandleFunc("POST /api/tasks", h.createTask)
	return mux
}

func (h *Handler) tasksJSON(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.tasks.List())
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if input.Text == "" {
		http.Error(w, "text is required", http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, h.tasks.Create(input.Text))
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}
