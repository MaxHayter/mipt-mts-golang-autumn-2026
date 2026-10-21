// Package httpapi adds multipart file upload to the JSON, XML and HTML routes.
package httpapi

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"http-seminar/internal/tasks"
)

const maxUploadSize = 5 << 20

type Handler struct {
	tasks     *tasks.Service
	page      *template.Template
	uploadDir string
}

type taskList struct {
	XMLName xml.Name     `xml:"tasks"`
	Tasks   []tasks.Task `xml:"task"`
}

func New(taskService *tasks.Service, page *template.Template, uploadDir string) *Handler {
	return &Handler{tasks: taskService, page: page, uploadDir: uploadDir}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", h.index)
	mux.HandleFunc("GET /api/tasks", h.tasksJSON)
	mux.HandleFunc("POST /api/tasks", h.createTask)
	mux.HandleFunc("GET /api/tasks.xml", h.tasksXML)
	mux.HandleFunc("POST /upload", h.upload)
	return mux
}

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	if err := h.page.Execute(w, h.tasks.List()); err != nil {
		http.Error(w, "cannot render page", http.StatusInternalServerError)
	}
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

func (h *Handler) upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "file is too large", http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "field file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := os.MkdirAll(h.uploadDir, 0o755); err != nil {
		http.Error(w, "cannot create upload directory", http.StatusInternalServerError)
		return
	}
	name := filepath.Base(header.Filename)
	destination, err := os.Create(filepath.Join(h.uploadDir, name))
	if err != nil {
		http.Error(w, "cannot save file", http.StatusInternalServerError)
		return
	}
	defer destination.Close()
	if _, err := io.Copy(destination, file); err != nil {
		http.Error(w, "cannot save file", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "uploaded: %s\n", name)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}
