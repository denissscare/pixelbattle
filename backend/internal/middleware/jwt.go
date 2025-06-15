package middleware

import (
	"context"
	"net/http"
	jwtutil "pixelbattle/pkg/jwt"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "userID"

func JWTAuth(jwtManager *jwtutil.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "missing or invalid auth header", http.StatusUnauthorized)
				return
			}
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwtManager.ParseToken(tokenString)
			if err != nil {
				http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
