package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/francisjdev/urlshortener/internal/http/handlers"
	"github.com/francisjdev/urlshortener/internal/repository"
	"github.com/francisjdev/urlshortener/internal/repository/memory"
	"github.com/francisjdev/urlshortener/internal/service"
)

func main() {
	//Create the memory repository
	var repo repository.URLRepository = memory.NewURLMemory()

	//Create the service, passing the repo
	svc := service.NewURLService(repo)

	//  Create the handler, passing the service
	handler := handlers.URLHandler{Service: svc}

	//Setup the HTTP server and mux
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.HealthHandler)
	mux.HandleFunc("PUT /create", handler.CreateURL)
	mux.HandleFunc("GET /{code}", handler.GetURL)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server listening on :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
