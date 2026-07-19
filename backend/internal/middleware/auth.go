package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWTSecret = []byte("super-secret-cloudpulse-key") // In production, load from env

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"error": "invalid auth header"}`, http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return JWTSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID := int(claims["user_id"].(float64))
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
	})
}

func GenerateJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString(JWTSecret)
}
