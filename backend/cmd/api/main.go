package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"cloudpulse/backend/internal/cache"
	"cloudpulse/backend/internal/handler"
	"cloudpulse/backend/internal/middleware"
	"cloudpulse/backend/internal/repository"
	"cloudpulse/backend/internal/storage"

	_ "github.com/lib/pq"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("CloudPulse API starting up...")

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8000"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/cloudpulse?sslmode=disable"
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		slog.Error("Error opening database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		slog.Error("Error connecting to database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	redisCache := cache.NewRedisCache(redisURL)
	storageSvc := storage.NewStorageService()

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
	router := handler.NewRouter(repo, redisCache, storageSvc)

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
