package main

import (
	"html/template"
	"log"
	"net/http"

	"http-seminar/internal/tasks"
	"http-seminar/steps/06_httptest/internal/httpapi"
)

func main() {
	page, err := template.ParseFiles("steps/06_httptest/web/templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	service := tasks.New([]tasks.Task{{ID: 1, Text: "Подготовить семинар"}})
	log.Fatal(http.ListenAndServe(":8080", httpapi.New(service, page, "steps/06_httptest/uploads").Routes()))
}
