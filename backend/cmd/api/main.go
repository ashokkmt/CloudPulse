package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"cloudpulse/backend/internal/handler"
	"cloudpulse/backend/internal/middleware"
	"cloudpulse/backend/internal/repository"

	_ "github.com/lib/pq"
)

func main() {
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8000"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/cloudpulse?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	tracingEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	shutdownTracing, err := middleware.InitTracing(context.Background(), "cloudpulse-backend", tracingEndpoint)
	if err != nil {
		log.Printf("Tracing disabled: %v", err)
	} else if shutdownTracing != nil {
		defer func() {
			if shutdownErr := shutdownTracing(context.Background()); shutdownErr != nil {
				log.Printf("Error shutting down tracing: %v", shutdownErr)
			}
		}()
	}

	repo := repository.NewPostgresRepository(db)
	router := handler.NewRouter(repo)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("CloudPulse API listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
