package middleware

import (
	"context"
	"net/http"
	jwtutil "pixelbattle/pkg/jwt"
)

type contextKey string

const UserIDKey contextKey = "userID"

func JWTAuth(jwtManager *jwtutil.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			tokenString := cookie.Value
			claims, err := jwtManager.ParseToken(tokenString)

			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
