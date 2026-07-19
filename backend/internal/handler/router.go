package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"cloudpulse/backend/internal/service"
)

func NewRouter(store *service.Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}
		jsonResponse(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/api/summary", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}
		jsonResponse(w, http.StatusOK, store.Summary())
	})
	mux.HandleFunc("/api/analytics", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}
		jsonResponse(w, http.StatusOK, store.Summary())
	})
	mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			jsonResponse(w, http.StatusOK, map[string]any{"tasks": store.ListTasks()})
		case http.MethodPost:
			var payload struct {
				Title string `json:"title"`
			}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || strings.TrimSpace(payload.Title) == "" {
				jsonError(w, http.StatusBadRequest, "title is required")
				return
			}
			created := store.AddTask(strings.TrimSpace(payload.Title))
			jsonResponse(w, http.StatusCreated, created)
		default:
			methodNotAllowed(w)
		}
	})
	mux.HandleFunc("/api/tasks/", func(w http.ResponseWriter, r *http.Request) {
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
			task, ok := store.UpdateTask(id, payload.Title, payload.Done)
			if !ok {
				jsonError(w, http.StatusNotFound, "task not found")
				return
			}
			jsonResponse(w, http.StatusOK, task)
		case http.MethodDelete:
			if !store.DeleteTask(id) {
				jsonError(w, http.StatusNotFound, "task not found")
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			methodNotAllowed(w)
		}
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}
		jsonResponse(w, http.StatusOK, map[string]string{
			"service": "CloudPulse API",
			"message": "Visit /healthz or /api/summary",
		})
	})

	return corsMiddleware(mux)
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
