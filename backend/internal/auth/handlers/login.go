package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	auth "pixelbattle/internal/auth/service"
	"pixelbattle/pkg/logger"
)

type LoginRequest struct {
	EmailOrUsername string `json:"email"`
	Password        string `json:"password"`
}

func LoginRender(svc *auth.Service, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := ""
		if tmpl, err := template.ParseFiles("static/signin.html"); err != nil {
			log.Errorf("HTTP GET /login: cannot open signin.html: %v", err)
			http.Error(w, "cannot open signin.html", http.StatusBadRequest)
		} else {
			tmpl.Execute(w, data)
		}
	}
}

func LoginHandler(svc *auth.Service, log *logger.Logger, minioHost string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Errorf("HTTP POST /login bad payload: %v", err)
			http.Error(w, "Bad payload", http.StatusBadRequest)
			return
		}

		if req.EmailOrUsername == "" || req.Password == "" {
			log.Infof("HTTP POST /login missing fields: %+v", req)
			http.Error(w, "Заполните все поля", http.StatusBadRequest)
			return
		}

		user, token, err := svc.Login(req.EmailOrUsername, req.Password)
		if err != nil {
			log.Infof("HTTP POST /login failed for user %s: %v", req.EmailOrUsername, err)
			http.Error(w, "Не удалось выполнить вход. Логин или пароль некорректны.", http.StatusUnauthorized)
			return
		}

		avatarURL := ""
		if user.AvatarURL != nil && *user.AvatarURL != "" {
			avatarURL = "http://" + minioHost + "/avatars/" + *user.AvatarURL
		}

		resp := map[string]interface{}{
			"token": token,
			"user": map[string]interface{}{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"avatar":   avatarURL,
			},
		}
		log.Infof("HTTP POST /login success for user %s (id=%d)", user.Username, user.ID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
