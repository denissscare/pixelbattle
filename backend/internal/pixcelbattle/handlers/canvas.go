package handlers

import (
	"encoding/json"
	"net/http"
	"pixelbattle/internal/pixcelbattle/service"
	"pixelbattle/pkg/logger"
)

func CanvasHandler(svc *service.BattleService, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		canvas, err := svc.InitCanvas(r.Context())
		if err != nil {
			log.Errorf("HTTP GET /canvas error: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(canvas)
		log.Infof("HTTP GET /canvas → returned %d pixels", len(canvas))
	}
}
