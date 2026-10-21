// Package httpapi adds XML as a second representation of the same tasks.
package httpapi

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"http-seminar/internal/tasks"
)

type Handler struct{ tasks *tasks.Service }

type taskList struct {
	XMLName xml.Name     `xml:"tasks"`
	Tasks   []tasks.Task `xml:"task"`
}

func New(taskService *tasks.Service) *Handler { return &Handler{tasks: taskService} }

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/tasks", h.tasksJSON)
	mux.HandleFunc("POST /api/tasks", h.createTask)
	mux.HandleFunc("GET /api/tasks.xml", h.tasksXML)
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

func (h *Handler) tasksXML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	xml.NewEncoder(w).Encode(taskList{Tasks: h.tasks.List()})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}
