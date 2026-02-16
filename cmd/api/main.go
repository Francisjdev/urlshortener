package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/francisjdev/urlshortener/internal/http/handlers"
	"github.com/francisjdev/urlshortener/internal/repository/postgres"
	"github.com/francisjdev/urlshortener/internal/service"
)

func main() {
	// Read database URL from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Dependency injection
	repo := postgres.NewPostgresURLRepository(db)
	svc := service.NewURLService(repo)
	handler := handlers.URLHandler{Service: svc}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/create", handler.CreateURL)
	mux.HandleFunc("/", handler.GetURL)

	// Read port from environment (Render sets this automatically)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Server listening on :%s\n", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
