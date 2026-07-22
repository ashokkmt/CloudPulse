package middleware

import (
	"fmt"
	"net/http"
	"time"

	"cloudpulse/backend/internal/cache"
)

// RateLimitMiddleware applies a rate limit of maxRequests per window.
// It uses Redis for distributed rate limiting based on the user's ID.
func RateLimitMiddleware(redisCache *cache.RedisCache, maxRequests int64, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract User ID from context (set by AuthMiddleware)
			userID, ok := r.Context().Value(UserIDKey).(int)
			if !ok {
				// If no user is authenticated, we might apply an IP-based limit, 
				// but since this is for protected routes, we can just deny or skip.
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Define the rate limit key
			key := fmt.Sprintf("ratelimit:user:%d:path:%s:method:%s", userID, r.URL.Path, r.Method)

			// Check rate limit
			allowed, err := redisCache.RateLimit(r.Context(), key, maxRequests, window)
			if err != nil {
				// In case of Redis error, we should probably fail open to not block traffic,
				// or fail closed depending on security requirements. We'll fail open and log.
				// For the demo, failing open is fine, but we should log the error.
				next.ServeHTTP(w, r)
				return
			}

			if !allowed {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error": "rate limit exceeded"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
