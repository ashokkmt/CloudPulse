package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"cloudpulse/backend/internal/middleware"
	"cloudpulse/backend/internal/repository"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/crypto/bcrypt"
)

func NewRouter(repo repository.Repository) http.Handler {
	mux := http.NewServeMux()

	// Public Routes
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}
		jsonResponse(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.Handle("/metrics", promhttp.Handler())

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
		// Dummy analytics for now, can be cached in Redis later
		userID := r.Context().Value(middleware.UserIDKey).(int)
		tasks, _ := repo.GetTasksByUserID(r.Context(), userID)
		completed := 0
		for _, t := range tasks {
			if t.Done {
				completed++
			}
		}
		jsonResponse(w, http.StatusOK, map[string]interface{}{
			"tasksTotal":     len(tasks),
			"tasksCompleted": completed,
			"tasksOpen":      len(tasks) - completed,
			"status":         "active",
		})
	})

	protectedMux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey).(int)
		
		switch r.Method {
		case http.MethodGet:
			tasks, err := repo.GetTasksByUserID(r.Context(), userID)
			if err != nil {
				jsonError(w, http.StatusInternalServerError, "error fetching tasks")
				return
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
			jsonResponse(w, http.StatusCreated, task)
		default:
			methodNotAllowed(w)
		}
	})

	protectedMux.HandleFunc("/api/tasks/", func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey).(int)
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/tasks/"))
		if err != nil {
			jsonError(w, http.StatusBadRequest, "invalid task id")
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
			jsonResponse(w, http.StatusOK, task)
		case http.MethodDelete:
			if err := repo.DeleteTask(r.Context(), id, userID); err != nil {
				jsonError(w, http.StatusNotFound, "task not found")
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			methodNotAllowed(w)
		}
	})

	// Wrap protected routes with Auth Middleware
	mux.Handle("/api/analytics", middleware.AuthMiddleware(protectedMux))
	mux.Handle("/api/tasks", middleware.AuthMiddleware(protectedMux))
	mux.Handle("/api/tasks/", middleware.AuthMiddleware(protectedMux))

	// Global Middlewares (Metrics, CORS)
	return corsMiddleware(middleware.MetricsMiddleware(mux))
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
