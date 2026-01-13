package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/francisjdev/urlshortener/internal/http/handlers"
)

func main() {
	fmt.Println("Server listening")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.HealthHandler)
	srv := http.Server{}
	srv.Addr = ":8080"
	srv.Handler = mux
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
