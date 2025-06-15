package handlers

import (
	"encoding/json"
	"net/http"
	auth "pixelbattle/internal/auth/service"
	"pixelbattle/pkg/logger"
)

type LoginRequest struct {
	EmailOrUsername string `json:"email"`
	Password        string `json:"password"`
}

func LoginHandler(svc *auth.Service, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Errorf("HTTP POST /login bad payload: %v", err)
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		if req.EmailOrUsername == "" || req.Password == "" {
			log.Infof("HTTP POST /login missing fields: %+v", req)
			http.Error(w, "missing fields", http.StatusBadRequest)
			return
		}
		user, token, err := svc.Login(req.EmailOrUsername, req.Password)
		if err != nil {
			log.Infof("HTTP POST /login failed for user %s: %v", req.EmailOrUsername, err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		resp := map[string]interface{}{
			"token": token,
			"user": map[string]interface{}{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"avatar":   user.AvatarURL,
			},
		}
		log.Infof("HTTP POST /login success for user %s (id=%d)", user.Username, user.ID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
