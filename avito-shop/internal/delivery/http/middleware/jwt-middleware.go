package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Нет токена", http.StatusUnauthorized)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Неверный формат токена", http.StatusUnauthorized)
				return
			}
			tokenStr := parts[1]
			claims := &jwt.RegisteredClaims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil || !token.Valid || claims.ExpiresAt.Time.Before(time.Now()) {
				http.Error(w, "Неверный или просроченный токен", http.StatusUnauthorized)
				return
			}
			id, err := strconv.ParseInt(claims.Subject, 10, 64)
			if err != nil {
				http.Error(w, "Некорректный идентификатор в токене", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "userID", id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
