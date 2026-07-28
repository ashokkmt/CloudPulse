package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cloudpulse/backend/internal/cache"
	"cloudpulse/backend/internal/middleware"
	"cloudpulse/backend/internal/model"
	"cloudpulse/backend/internal/repository"
	"cloudpulse/backend/internal/storage"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/crypto/bcrypt"
)

func NewRouter(repo repository.Repository, redisCache *cache.RedisCache, storageSvc *storage.StorageService, db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	// Public Routes — Health Check (used by K8s liveness/readiness probes)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		// Check database connectivity
		if err := db.PingContext(ctx); err != nil {
			slog.Error("Health check failed: database unreachable", slog.String("error", err.Error()))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status": "unhealthy",
				"error":  "database unreachable: " + err.Error(),
			})
			return
		}

		// Check Redis connectivity
		if err := redisCache.Ping(ctx); err != nil {
			slog.Error("Health check failed: redis unreachable", slog.String("error", err.Error()))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status": "unhealthy",
				"error":  "redis unreachable: " + err.Error(),
			})
			return
		}

		jsonResponse(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/api/broken", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Simulated Server Error", http.StatusInternalServerError)
	})

	mux.HandleFunc("/api/auth/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			methodNotAllowed(w)
			return
		}
		var payload struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Email == "" || payload.Password == "" {
			jsonError(w, http.StatusBadRequest, "invalid request")
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		if err != nil {
			jsonError(w, http.StatusInternalServerError, "error hashing password")
			return
		}

		user, err := repo.CreateUser(r.Context(), payload.Email, string(hash))
		if err != nil {
			jsonError(w, http.StatusConflict, "user already exists")
			return
		}

		jsonResponse(w, http.StatusCreated, user)
	})

	mux.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			methodNotAllowed(w)
			return
		}
		var payload struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			jsonError(w, http.StatusBadRequest, "invalid request")
			return
		}

		user, err := repo.GetUserByEmail(r.Context(), payload.Email)
		if err != nil {
			jsonError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
			jsonError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}

		token, err := middleware.GenerateJWT(user.ID)
		if err != nil {
			jsonError(w, http.StatusInternalServerError, "error generating token")
			return
		}

		jsonResponse(w, http.StatusOK, map[string]interface{}{
			"token": token,
			"user":  user,
		})
	})

	// Protected Routes
	protectedMux := http.NewServeMux()

	protectedMux.HandleFunc("/api/analytics", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}

		userID := r.Context().Value(middleware.UserIDKey).(int)
		cacheKey := fmt.Sprintf("dashboard:stats:user:%d", userID)

		var stats map[string]interface{}
		err := redisCache.Get(r.Context(), cacheKey, &stats)
		if err == nil {
			// Cache hit
			jsonResponse(w, http.StatusOK, stats)
			return
		}

		// Cache miss
		tasks, _ := repo.GetTasksByUserID(r.Context(), userID)
		completed := 0
		for _, t := range tasks {
			if t.Done {
				completed++
			}
		}
		stats = map[string]interface{}{
			"tasksTotal":     len(tasks),
			"tasksCompleted": completed,
			"tasksOpen":      len(tasks) - completed,
			"status":         "active",
		}

		// Store in Redis with 30s TTL
		_ = redisCache.Set(r.Context(), cacheKey, stats, 30*time.Second)

		jsonResponse(w, http.StatusOK, stats)
	})

	protectedMux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey).(int)

		switch r.Method {
		case http.MethodGet:
			cacheKey := fmt.Sprintf("tasks:user:%d", userID)
			var tasks []model.Task
			err := redisCache.Get(r.Context(), cacheKey, &tasks)
			if err == nil {
				// Inject presigned URLs for thumbnails before returning
				for i := range tasks {
					if tasks[i].ThumbnailURL != nil && *tasks[i].ThumbnailURL != "" {
						url, _ := storageSvc.GeneratePresignedURL(r.Context(), *tasks[i].ThumbnailURL)
						tasks[i].ThumbnailURL = &url
					}
				}
				jsonResponse(w, http.StatusOK, map[string]any{"tasks": tasks})
				return
			}

			tasks, err = repo.GetTasksByUserID(r.Context(), userID)
			if err != nil {
				jsonError(w, http.StatusInternalServerError, "error fetching tasks")
				return
			}

			_ = redisCache.Set(r.Context(), cacheKey, tasks, 30*time.Second)

			// Inject presigned URLs for thumbnails before returning
			for i := range tasks {
				if tasks[i].ThumbnailURL != nil && *tasks[i].ThumbnailURL != "" {
					url, _ := storageSvc.GeneratePresignedURL(r.Context(), *tasks[i].ThumbnailURL)
					tasks[i].ThumbnailURL = &url
				}
			}

			jsonResponse(w, http.StatusOK, map[string]any{"tasks": tasks})
		case http.MethodPost:
			var payload struct {
				Title string `json:"title"`
			}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || strings.TrimSpace(payload.Title) == "" {
				jsonError(w, http.StatusBadRequest, "title is required")
				return
			}
			task, err := repo.CreateTask(r.Context(), userID, strings.TrimSpace(payload.Title))
			if err != nil {
				jsonError(w, http.StatusInternalServerError, "error creating task")
				return
			}

			// Invalidate caches
			_ = redisCache.Delete(r.Context(), fmt.Sprintf("dashboard:stats:user:%d", userID))
			_ = redisCache.Delete(r.Context(), fmt.Sprintf("tasks:user:%d", userID))

			// Emit job to Redis Queue
			job := model.TaskJob{
				TaskID:    task.ID,
				UserID:    task.UserID,
				CreatedAt: time.Now(),
			}
			_ = redisCache.Enqueue(r.Context(), "tasks_queue", job)

			slog.Info("Task Created", slog.Int("task_id", task.ID), slog.Int("user_id", task.UserID))
			jsonResponse(w, http.StatusCreated, task)
		default:
			methodNotAllowed(w)
		}
	})

	protectedMux.HandleFunc("/api/tasks/", func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey).(int)
		path := strings.TrimPrefix(r.URL.Path, "/api/tasks/")

		parts := strings.Split(path, "/")
		idStr := parts[0]
		action := ""
		if len(parts) > 1 {
			action = strings.Join(parts[1:], "/")
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			jsonError(w, http.StatusBadRequest, "invalid task id")
			return
		}

		if action == "attachment/url" {
			if r.Method != http.MethodGet {
				methodNotAllowed(w)
				return
			}
			task, err := repo.GetTaskByID(r.Context(), id, userID)
			if err != nil {
				jsonError(w, http.StatusNotFound, "task not found")
				return
			}
			if task.AttachmentName == nil || *task.AttachmentName == "" {
				jsonError(w, http.StatusNotFound, "attachment not found")
				return
			}
			url, err := storageSvc.GeneratePresignedURL(r.Context(), *task.AttachmentName)
			if err != nil {
				slog.Error("GeneratePresignedURL failed",
					slog.String("object", *task.AttachmentName),
					slog.String("bucket", storageSvc.BucketName()),
					slog.Any("err", err),
				)
				jsonError(w, http.StatusInternalServerError, "error generating url")
				return
			}
			jsonResponse(w, http.StatusOK, map[string]string{"url": url})
			return
		}

		if action == "attachment" {
			switch r.Method {
			case http.MethodPost:
				r.Body = http.MaxBytesReader(w, r.Body, 20<<20) // 20 MB limit
				if err := r.ParseMultipartForm(20 << 20); err != nil {
					jsonError(w, http.StatusBadRequest, "file too large or invalid")
					return
				}

				file, header, err := r.FormFile("attachment")
				if err != nil {
					jsonError(w, http.StatusBadRequest, "missing attachment")
					return
				}
				defer file.Close()

				// verify task exists and belongs to user
				_, err = repo.GetTaskByID(r.Context(), id, userID)
				if err != nil {
					jsonError(w, http.StatusNotFound, "task not found")
					return
				}

				objectName := fmt.Sprintf("user-%d/task-%d/%s", userID, id, header.Filename)
				mimeType := header.Header.Get("Content-Type")
				size := header.Size

				_, err = storageSvc.UploadFile(r.Context(), objectName, file, size, mimeType)
				if err != nil {
					jsonError(w, http.StatusInternalServerError, "error uploading file")
					return
				}

				task, err := repo.UpdateTaskAttachment(r.Context(), id, userID, nil, &objectName, &size, &mimeType)
				if err != nil {
					jsonError(w, http.StatusInternalServerError, "error updating task")
					return
				}

				_ = redisCache.Delete(r.Context(), fmt.Sprintf("dashboard:stats:user:%d", userID))
				_ = redisCache.Delete(r.Context(), fmt.Sprintf("tasks:user:%d", userID))

				// Push thumbnail job if image
				if strings.HasPrefix(mimeType, "image/") {
					job := model.TaskJob{
						TaskID:    task.ID,
						UserID:    task.UserID,
						CreatedAt: time.Now(),
					}
					_ = redisCache.Enqueue(r.Context(), "thumbnail_queue", job)
				}

				jsonResponse(w, http.StatusOK, task)
			case http.MethodDelete:
				task, err := repo.GetTaskByID(r.Context(), id, userID)
				if err != nil {
					jsonError(w, http.StatusNotFound, "task not found")
					return
				}

				if task.AttachmentName != nil && *task.AttachmentName != "" {
					_ = storageSvc.DeleteFile(r.Context(), *task.AttachmentName)
				}

				task, err = repo.UpdateTaskAttachment(r.Context(), id, userID, nil, nil, nil, nil)
				if err != nil {
					jsonError(w, http.StatusInternalServerError, "error updating task")
					return
				}

				_ = redisCache.Delete(r.Context(), fmt.Sprintf("dashboard:stats:user:%d", userID))
				_ = redisCache.Delete(r.Context(), fmt.Sprintf("tasks:user:%d", userID))

				jsonResponse(w, http.StatusOK, task)
			default:
				methodNotAllowed(w)
			}
			return
		}

		switch r.Method {
		case http.MethodPut:
			var payload struct {
				Title *string `json:"title"`
				Done  *bool   `json:"done"`
			}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				jsonError(w, http.StatusBadRequest, "invalid request body")
				return
			}
			task, err := repo.UpdateTask(r.Context(), id, userID, payload.Title, payload.Done)
			if err != nil {
				jsonError(w, http.StatusNotFound, "task not found")
				return
			}

			// Invalidate caches
			_ = redisCache.Delete(r.Context(), fmt.Sprintf("dashboard:stats:user:%d", userID))
			_ = redisCache.Delete(r.Context(), fmt.Sprintf("tasks:user:%d", userID))

			slog.Info("Task Updated", slog.Int("task_id", task.ID), slog.Int("user_id", task.UserID))
			jsonResponse(w, http.StatusOK, task)
		case http.MethodDelete:
			// Fetch the task first to check if there is an attachment
			task, err := repo.GetTaskByID(r.Context(), id, userID)
			if err != nil {
				jsonError(w, http.StatusNotFound, "task not found")
				return
			}

			// If there's an attachment, delete it from spaces
			if task.AttachmentName != nil && *task.AttachmentName != "" {
				_ = storageSvc.DeleteFile(r.Context(), *task.AttachmentName)
			}

			if err := repo.DeleteTask(r.Context(), id, userID); err != nil {
				jsonError(w, http.StatusNotFound, "task not found")
				return
			}

			// Invalidate caches
			_ = redisCache.Delete(r.Context(), fmt.Sprintf("dashboard:stats:user:%d", userID))
			_ = redisCache.Delete(r.Context(), fmt.Sprintf("tasks:user:%d", userID))

			slog.Info("Task Deleted", slog.Int("task_id", id), slog.Int("user_id", userID))
			w.WriteHeader(http.StatusNoContent)
		default:
			methodNotAllowed(w)
		}
	})

	// Rate limiting middleware for write APIs
	writeRateLimiter := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
				middleware.RateLimitMiddleware(redisCache, 60, time.Minute)(next).ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	// Wrap protected routes with Auth Middleware and Rate Limiter
	mux.Handle("/api/analytics", middleware.AuthMiddleware(protectedMux))
	mux.Handle("/api/tasks", middleware.AuthMiddleware(writeRateLimiter(protectedMux)))
	mux.Handle("/api/tasks/", middleware.AuthMiddleware(writeRateLimiter(protectedMux)))

	// Global Middlewares (Metrics, CORS)
	return corsMiddleware(middleware.TracingMiddleware(middleware.MetricsMiddleware(mux)))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "600")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func jsonResponse(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func jsonError(w http.ResponseWriter, status int, message string) {
	jsonResponse(w, status, map[string]string{"error": message})
}

func methodNotAllowed(w http.ResponseWriter) {
	jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
}
