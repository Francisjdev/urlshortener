package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/francisjdev/urlshortener/internal/http/handlers"
	"github.com/francisjdev/urlshortener/internal/repository/postgres"
	"github.com/francisjdev/urlshortener/internal/service"
)

func main() {
	// Connect to Postgres
	db, err := sql.Open(
		"pgx",
		"postgres://localhost:5432/urlshortener?sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Create the Postgres repository
	repo := postgres.NewPostgresURLRepository(db)

	// Create the service, passing the repo
	svc := service.NewURLService(repo)

	// Create the handler, passing the service
	handler := handlers.URLHandler{Service: svc}

	// Setup HTTP server and mux
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/create", handler.CreateURL)
	mux.HandleFunc("/", handler.GetURL) // capture all GETs for short codes

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server listening on :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
