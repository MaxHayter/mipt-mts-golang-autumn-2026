package main

import (
	"log"
	"net/http"

	"http-seminar/internal/tasks"
	"http-seminar/steps/03_xml/internal/httpapi"
)

func main() {
	service := tasks.New([]tasks.Task{{ID: 1, Text: "Подготовить семинар"}})
	log.Fatal(http.ListenAndServe(":8080", httpapi.New(service).Routes()))
}
