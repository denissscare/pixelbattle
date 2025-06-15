package handlers

import (
	"encoding/json"
	"net/http"
	auth "pixelbattle/internal/auth/service"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		if req.Username == "" || req.Email == "" || req.Password == "" {
			http.Error(w, "missing fields", http.StatusBadRequest)
			return
		}
		if err := svc.Register(req.Username, req.Email, req.Password); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
