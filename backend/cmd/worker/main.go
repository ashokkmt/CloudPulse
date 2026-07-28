package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"image/jpeg"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"time"

	"cloudpulse/backend/internal/cache"
	"cloudpulse/backend/internal/metrics"
	"cloudpulse/backend/internal/middleware"
	"cloudpulse/backend/internal/model"
	"cloudpulse/backend/internal/repository"
	"cloudpulse/backend/internal/storage"

	_ "github.com/lib/pq"
	"golang.org/x/image/draw"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("CloudPulse Background Worker started")

	tracingEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	shutdownTracing, err := middleware.InitTracing(context.Background(), "cloudpulse-worker", tracingEndpoint)
	if err != nil {
		slog.Warn("Tracing disabled", slog.String("error", err.Error()))
	} else if shutdownTracing != nil {
		defer func() {
			if shutdownErr := shutdownTracing(context.Background()); shutdownErr != nil {
				slog.Error("Error shutting down tracing", slog.String("error", shutdownErr.Error()))
			}
		}()
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/cloudpulse?sslmode=disable"
	}

	// Log config (mask sensitive data)
	maskedDBHost := "unknown"
	if parsed, parseErr := url.Parse(dbURL); parseErr == nil {
		maskedDBHost = parsed.Host
	}
	slog.Info("Worker configuration loaded",
		"dbHost", maskedDBHost,
		"redisURL_len", len(redisURL),
	)

	// --- Redis Connection ---
	slog.Info("Worker connecting to Redis...", "url_length", len(redisURL))
	redisCache := cache.NewRedisCache(redisURL)
	slog.Info("Worker Redis connection established")

	// --- Storage Service ---
	slog.Info("Worker initializing storage service...")
	storageSvc := storage.NewStorageService()
	slog.Info("Worker storage service initialized")

	// --- Database Connection ---
	slog.Info("Worker connecting to PostgreSQL...", "host", maskedDBHost)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		slog.Error("FATAL: Worker failed to open database connection pool", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("Worker PostgreSQL connection pool opened, pinging...")
	if err := db.Ping(); err != nil {
		slog.Error("FATAL: Worker failed to ping database", slog.String("error", err.Error()), slog.String("host", maskedDBHost))
		os.Exit(1)
	}
	slog.Info("Worker PostgreSQL connected successfully")
	repo := repository.NewPostgresRepository(db)
	ctx := context.Background()

	// Start task processing loop
	go func() {
		for {
			jobStr, err := redisCache.Dequeue(ctx, "tasks_queue", 0)
			if err != nil {
				slog.Error("Error dequeueing tasks_queue job", slog.String("error", err.Error()))
				continue
			}
			start := time.Now()
			var job model.TaskJob
			if err := json.Unmarshal([]byte(jobStr), &job); err != nil {
				slog.Error("Error parsing tasks_queue job", slog.String("error", err.Error()))
				metrics.WorkerJobsProcessed.WithLabelValues("tasks_queue", "failure").Inc()
				continue
			}
			slog.Info("Worker Processing", slog.String("queue", "tasks_queue"), slog.Int("task_id", job.TaskID), slog.Int("user_id", job.UserID))
			// Simulating work
			metrics.WorkerJobsProcessed.WithLabelValues("tasks_queue", "success").Inc()
			metrics.QueueProcessingDuration.WithLabelValues("tasks_queue").Observe(time.Since(start).Seconds())
			slog.Info("Worker Success", slog.String("queue", "tasks_queue"), slog.Int("task_id", job.TaskID))
		}
	}()

	// Start thumbnail processing loop
	for {
		jobStr, err := redisCache.Dequeue(ctx, "thumbnail_queue", 0)
		if err != nil {
			slog.Error("Error dequeueing thumbnail_queue job", slog.String("error", err.Error()))
			continue
		}
		start := time.Now()
		var job model.TaskJob
		if err := json.Unmarshal([]byte(jobStr), &job); err != nil {
			slog.Error("Error parsing thumbnail_queue job", slog.String("error", err.Error()))
			metrics.WorkerJobsProcessed.WithLabelValues("thumbnail_queue", "failure").Inc()
			continue
		}
		
		slog.Info("Worker Processing", slog.String("queue", "thumbnail_queue"), slog.Int("task_id", job.TaskID))
		err = processThumbnail(ctx, job, repo, storageSvc, redisCache)
		if err != nil {
			slog.Error("Worker Failure", slog.String("queue", "thumbnail_queue"), slog.Int("task_id", job.TaskID), slog.String("error", err.Error()))
			metrics.WorkerJobsProcessed.WithLabelValues("thumbnail_queue", "failure").Inc()
		} else {
			slog.Info("Worker Success", slog.String("queue", "thumbnail_queue"), slog.Int("task_id", job.TaskID))
			metrics.WorkerJobsProcessed.WithLabelValues("thumbnail_queue", "success").Inc()
		}
		metrics.QueueProcessingDuration.WithLabelValues("thumbnail_queue").Observe(time.Since(start).Seconds())
	}
}

func processThumbnail(ctx context.Context, job model.TaskJob, repo repository.Repository, storageSvc *storage.StorageService, redisCache *cache.RedisCache) error {
	task, err := repo.GetTaskByID(ctx, job.TaskID, job.UserID)
	if err != nil {
		return fmt.Errorf("could not find task: %w", err)
	}

	if task.AttachmentName == nil || *task.AttachmentName == "" {
		return fmt.Errorf("no attachment found")
	}

	file, err := storageSvc.GetFile(ctx, *task.AttachmentName)
	if err != nil {
		return fmt.Errorf("could not download file: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("could not decode image: %w", err)
	}

	// Calculate thumbnail size (max 200x200)
	bounds := img.Bounds()
	ratio := float64(bounds.Dx()) / float64(bounds.Dy())
	newW, newH := 200, 200
	if ratio > 1 {
		newH = int(200 / ratio)
	} else {
		newW = int(200 * ratio)
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.BiLinear.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 80}); err != nil {
		return fmt.Errorf("could not encode jpeg: %w", err)
	}

	thumbName := strings.TrimSuffix(*task.AttachmentName, ".jpg") // simplistic, might need better ext handling
	thumbName = strings.TrimSuffix(thumbName, ".png")
	thumbName = strings.TrimSuffix(thumbName, ".jpeg")
	thumbName = thumbName + "-thumb.jpg"

	thumbKey, err := storageSvc.UploadFile(ctx, thumbName, &buf, int64(buf.Len()), "image/jpeg")
	if err != nil {
		return fmt.Errorf("could not upload thumbnail: %w", err)
	}

	if err := repo.UpdateTaskThumbnail(ctx, task.ID, task.UserID, &thumbKey); err != nil {
		return fmt.Errorf("could not update task thumbnail url: %w", err)
	}

	// Invalidate cache
	_ = redisCache.Delete(ctx, fmt.Sprintf("tasks:user:%d", task.UserID))
	return nil
}
