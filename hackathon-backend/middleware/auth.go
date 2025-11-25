package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// "Bearer xxxxx.yyyyy.zzzzz"
		var tokenString string
		_, err := fmt.Sscanf(authHeader, "Bearer %s", &tokenString)
		if err != nil {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		// JWT secret
		secret := []byte(os.Getenv("JWT_SECRET"))

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := claims["user_id"].(string)

		// Context に userID を入れる
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
