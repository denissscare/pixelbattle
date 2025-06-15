package handlers

import (
	"encoding/json"
	"net/http"
	auth "pixelbattle/internal/auth/service"
	"pixelbattle/pkg/logger"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(svc *auth.Service, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Errorf("HTTP POST /register bad payload: %v", err)
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		if req.Username == "" || req.Email == "" || req.Password == "" {
			log.Infof("HTTP POST /register missing fields: %+v", req)
			http.Error(w, "missing fields", http.StatusBadRequest)
			return
		}
		if err := svc.Register(req.Username, req.Email, req.Password); err != nil {
			log.Infof("HTTP POST /register failed: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		log.Infof("HTTP POST /register → registered user %s (%s)", req.Username, req.Email)
	}
}
