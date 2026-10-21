package main

import (
	"flag"
	"html/template"
	"log"
	"net"
	"net/http"

	"http-seminar/internal/tasks"
	"http-seminar/steps/07_tls/internal/httpapi"
)

func main() {
	tls := flag.Bool("tls", false, "serve HTTPS with cert.pem and key.pem")
	flag.Parse()

	page, err := template.ParseFiles("steps/07_tls/web/templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	service := tasks.New([]tasks.Task{{ID: 1, Text: "Подготовить семинар"}})
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	server := &http.Server{Handler: httpapi.New(service, page, "steps/07_tls/uploads").Routes()}
	if *tls {
		log.Fatal(server.ServeTLS(listener, "cert.pem", "key.pem"))
	}
	log.Fatal(server.Serve(listener))
}
