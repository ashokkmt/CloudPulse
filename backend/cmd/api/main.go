package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"net/url"
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

	// Log config (mask sensitive data)
	maskedDBHost := "unknown"
	if parsed, err := url.Parse(dbURL); err == nil {
		maskedDBHost = parsed.Host
	}
	slog.Info("Configuration loaded",
		"port", port,
		"dbHost", maskedDBHost,
		"redisURL_len", len(redisURL),
	)

	// --- Database Connection ---
	slog.Info("Connecting to PostgreSQL...", "host", maskedDBHost)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		slog.Error("FATAL: Failed to open database connection pool", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("PostgreSQL connection pool opened, pinging...")
	if err := db.Ping(); err != nil {
		slog.Error("FATAL: Failed to ping database — is the DB reachable from this network?", slog.String("error", err.Error()), slog.String("host", maskedDBHost))
		os.Exit(1)
	}
	slog.Info("PostgreSQL connected successfully")

	// --- Redis Connection ---
	slog.Info("Connecting to Redis...", "url_length", len(redisURL))
	redisCache := cache.NewRedisCache(redisURL)
	slog.Info("Redis connection established")

	// --- Storage Service ---
	slog.Info("Initializing storage service...")
	storageSvc := storage.NewStorageService()
	slog.Info("Storage service initialized")

	slog.Info("All services initialized, starting HTTP server...")

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
	router := handler.NewRouter(repo, redisCache, storageSvc, db)

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
