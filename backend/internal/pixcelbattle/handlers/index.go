package handlers

import (
	"net/http"
	"pixelbattle/internal/pixcelbattle/service"
	"pixelbattle/pkg/logger"

	"html/template"
)

func CanvasRender(svc *service.BattleService, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := ""

		if tmpl, err := template.ParseFiles("static/canvas.html"); err != nil {
			log.Errorf("HTTP GET /register: cannot open canvas.html: %v", err)
			http.Error(w, "cannot open canvas.html", http.StatusBadRequest)
		} else {
			tmpl.Execute(w, data)
		}
	}
}
