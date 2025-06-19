package handlers

import (
	"fmt"
	"net/http"
	auth "pixelbattle/internal/auth/service"
	"pixelbattle/internal/s3"
	"pixelbattle/pkg/logger"

	"html/template"
)

func RegisterRender(svc *auth.Service, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := ""

		if tmpl, err := template.ParseFiles("static/signup.html"); err != nil {
			log.Errorf("HTTP GET /register: cannot open signup.html: %v", err)
			http.Error(w, "cannot open signup.html", http.StatusBadRequest)
		} else {
			tmpl.Execute(w, data)
		}
	}
}

func RegisterHandler(s3Client *s3.Client, svc *auth.Service, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			log.Errorf("HTTP POST /register: cannot parse multipart form: %v", err)
			http.Error(w, "Не удалось загрузить форму", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		if username == "" || email == "" || password == "" {
			log.Infof("HTTP POST /register: missing fields: username=%s, email=%s", username, email)
			http.Error(w, "Заполнены не все поля", http.StatusBadRequest)
			return
		}

		userID, err := svc.RegisterWithID(username, email, password)
		if err != nil {
			log.Infof("HTTP POST /register: failed: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		file, fileHeader, err := r.FormFile("avatar")
		if err == nil {

			defer file.Close()
			objectName := fmt.Sprintf("%d_%s", userID, fileHeader.Filename)
			_, err = s3Client.UploadFile(r.Context(), fileHeader, objectName)

			if err != nil {
				log.Errorf("avatar upload: s3 error: %v", err)
				http.Error(w, "Ошибка загрузки", http.StatusInternalServerError)
				return
			}
			if err := svc.UpdateAvatarURL(userID, objectName); err != nil {
				log.Errorf("avatar upload: update db error: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			log.Infof("avatar uploaded: userID=%d, file=%s", userID, fileHeader.Filename)
		}
		log.Infof("HTTP POST /register → registered user %s (%s), id=%d", username, email, userID)
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}
}
