package handlers

import (
	"encoding/json"
	"net/http"
	auth "pixelbattle/internal/auth/service"
	"pixelbattle/internal/middleware"
	"pixelbattle/pkg/logger"
)

type updateEmailRequest struct {
	Email string `json:"email"`
}

func UpdateEmailHandler(authSvc *auth.Service, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDVal := r.Context().Value(middleware.UserIDKey)
		if userIDVal == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		userID, ok := userIDVal.(int)
		if !ok {
			http.Error(w, "invalid user id", http.StatusUnauthorized)
			return
		}

		var req updateEmailRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if req.Email == "" {
			http.Error(w, "empty email", http.StatusBadRequest)
			return
		}
		if err := authSvc.UpdateEmail(userID, req.Email); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
