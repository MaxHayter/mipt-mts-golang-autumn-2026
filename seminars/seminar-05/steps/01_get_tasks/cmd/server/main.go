package main

import (
	"log"
	"net/http"

	"http-seminar/internal/tasks"
	"http-seminar/steps/01_get_tasks/internal/httpapi"
)

func main() {
	service := tasks.New([]tasks.Task{{ID: 1, Text: "Подготовить семинар"}})
	log.Fatal(http.ListenAndServe(":8080", httpapi.New(service).Routes()))
}
