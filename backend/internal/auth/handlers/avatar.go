package handlers

import (
	"net/http"
	auth "pixelbattle/internal/auth/service"
	"pixelbattle/internal/middleware"
	"pixelbattle/pkg/logger"
)

func UploadAvatarHandler(authSvc *auth.Service, log *logger.Logger) http.HandlerFunc {
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

		_, fileHeader, err := r.FormFile("avatar")
		if err != nil {
			log.Errorf("avatar upload: bad file: %v", err)
			http.Error(w, "bad file", http.StatusBadRequest)
			return
		}

		if err := authSvc.UploadAvatar(r.Context(), userID, fileHeader); err != nil {
			log.Errorf("avatar upload: %v", err)
			http.Error(w, "upload error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}